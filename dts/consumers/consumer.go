package consumers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/murderxchip/aliyun-dts-dispatcher/config"
	"github.com/murderxchip/aliyun-dts-dispatcher/define"
	"strconv"
	"strings"

	dtsavro "github.com/murderxchip/aliyun-dts-avro/avro"
	"github.com/murderxchip/aliyun-dts-dispatcher/log"
	"github.com/Shopify/sarama"
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	cluster "github.com/bsm/sarama-cluster"
	"io"
	"sync"
	"time"
)

const (
	STATUS_RUNNING = iota << 1
	STATUS_STOPPED
)

type Consumer struct {
	Name     string
	Config   *ConsumerConfig
	Consumer *cluster.Consumer
	stop     chan int
	exit     chan interface{}
	mux      sync.Mutex
	queue    chan<- DtsData
}

type DtsData struct {
	//数据源
	Consumer string `json:"consumer"`
	//表名
	Table string `json:"table"`
	//事件时间
	Datetime string `json:"datetime"`
	//字段信息
	Data []DtsValue `json:"data"`
	//变更字段列表
	Fields []string `json:"fields"`
	//数据库事件
	Operation string `json:"operation"`
}

type DtsValue struct {
	Name      string `json:"name"`
	BeforeVal string `json:"before"`
	AfterVal  string `json:"after"`
}

func NewConsumer(config ConsumerConfig, q chan<- DtsData) (consumer *Consumer, err error) {
	consumer = &Consumer{}
	consumer.Config = &config
	consumer.stop = make(chan int, 1)
	consumer.exit = make(chan interface{}, 1)
	consumer.queue = q
	consumer.Name = config.Name

	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Group.Return.Notifications = true
	cfg.Net.MaxOpenRequests = 100
	//cfg.Consumer.Offsets.CommitInterval = 1 * time.Second
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = time.Second * 5

	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Net.SASL.Enable = true
	cfg.Net.SASL.User = config.GetSASLUser()
	cfg.Net.SASL.Password = config.Password
	cfg.Version = sarama.V0_11_0_0

	consumer.Consumer, err = cluster.NewConsumer(config.Broker, config.GroupId, config.Topic, cfg)
	if err != nil {
		log.Fatal(define.ConsumerTag, "consumer failed ")
		panic(err)
	}

	log.Info(define.ConsumerTag, "starting consumer ", consumer.Name)
	return
}

func (c *Consumer) Start() {

	go func() {
		for err := range c.Consumer.Errors() {
			//todo
			log.Fatal(define.ConsumerTag,c.Consumer)
			panic(err)
		}
	}()

	//ticker := time.NewTicker(time.Second * 5)
	//go func() {
	//	for { //循环
	//		select {
	//		case <-ticker.C:
	//			{
	//				c.Consumer.CommitOffsets()
	//			}
	//		}
	//
	//	}
	//}()
	// consume notifications
	go func() {
		for ntf := range c.Consumer.Notifications() {
			log.Info(define.ConsumerTag, "Rebalanced:", ntf)
		}
	}()

	go func() {
		// Pre compile schema of avro
		t := dtsavro.NewRecord()
		deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
		if err != nil {
			log.Fatal(define.ConsumerTag, deser, err)
			panic(err)
		}
		// consume messages, watch signals
		var r io.Reader
		for {
			select {
			case msg, ok := <-c.Consumer.Messages():
				if ok {
					r = bytes.NewReader(msg.Value)
					t = dtsavro.NewRecord()
					vm.Eval(r, deser, t)
					switch t.Operation {
					case dtsavro.OperationHEARTBEAT:
					case dtsavro.OperationDELETE:
					case dtsavro.OperationUPDATE, dtsavro.OperationINSERT:

						dbname := strings.Split(t.ObjectName.String, ".")

						dtsData := DtsData{}
						dtsData.Consumer = c.Name
						dtsData.Table = dbname[1]
						dtsData.Operation = GetOperationText(t.Operation)
						dtsData.Datetime = msg.Timestamp.Format("2006-01-02 15:04:05")

						for k, j := range t.Fields.ArrayField {
							var beforeVal, afterVal *dtsavro.UnionNullIntegerCharacterDecimalFloatTimestampDateTimeTimestampWithTimeZoneBinaryGeometryTextGeometryBinaryObjectTextObjectEmptyObject
							var before, after dtsavro.UnionNullIntegerCharacterDecimalFloatTimestampDateTimeTimestampWithTimeZoneBinaryGeometryTextGeometryBinaryObjectTextObjectEmptyObject

							if t.Operation == dtsavro.OperationUPDATE {
								beforeVal = t.BeforeImages.ArrayUnionNullIntegerCharacterDecimalFloatTimestampDateTimeTimestampWithTimeZoneBinaryGeometryTextGeometryBinaryObjectTextObjectEmptyObject[k]
								a, _ := json.Marshal(beforeVal)
								_ = json.Unmarshal(a, &before)
							}
							afterVal = t.AfterImages.ArrayUnionNullIntegerCharacterDecimalFloatTimestampDateTimeTimestampWithTimeZoneBinaryGeometryTextGeometryBinaryObjectTextObjectEmptyObject[k]

							b, _ := json.Marshal(afterVal)
							_ = json.Unmarshal(b, &after)

							var v DtsValue
							v.Name = j.Name
							if t.Operation == dtsavro.OperationUPDATE {
								v.BeforeVal = Extract(before)
							}
							v.AfterVal = Extract(after)

							if v.BeforeVal != v.AfterVal {
								dtsData.Fields = append(dtsData.Fields, v.Name)
							}
							dtsData.Data = append(dtsData.Data, v)
						}
						/*
						   id 3
						   name 253
						   age 3
						   column_4 12
						   column_5 4
						   column_6 252
						   column_7 7
						   column_8 15
						*/
						//todo
						log.Info(define.ConsumerTag, "push to queue:", dtsData)
						c.queue <- dtsData
						if config.Config().Env == "dev" {
							time.Sleep(time.Second * 20)
						}
					}
				}
			case <-c.stop:
				c.Consumer.CommitOffsets()
				c.Consumer.Close()
				c.exit <- 1
				return
			}
		}
	}()
}

func (c *Consumer) GetName() string {
	return c.Name
}

func (c *Consumer) Shutdown() {
	c.stop <- 1
	log.Info(define.ConsumerTag, fmt.Sprintf("consumer %s exit", c.Name))
	<-c.exit
}

//func (c *Consumer) WaitExit() {
//	tools.WaitTimeout(fmt.Sprintf("%s: %s", define.ConsumerTag, c.Name), c.exit, 5*time.Second)
//}

func GetOperationText(op dtsavro.Operation) string {
	switch op {
	case dtsavro.OperationINSERT:
		return "INSERT"
	case dtsavro.OperationUPDATE:
		return "UPDATE"
	case dtsavro.OperationDELETE:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func Extract(v dtsavro.UnionNullIntegerCharacterDecimalFloatTimestampDateTimeTimestampWithTimeZoneBinaryGeometryTextGeometryBinaryObjectTextObjectEmptyObject) string {
	if v.Integer != nil {
		return v.Integer.Value
	}

	if v.Float != nil {
		return strconv.FormatFloat(v.Float.Value, 'E', -1, 32)
	}

	if v.TextObject != nil {
		return v.TextObject.Value
	}

	if v.Character != nil {
		return string(v.Character.Value)
	}

	if v.Timestamp != nil {
		return strconv.FormatInt(v.Timestamp.Timestamp, 10)
	}

	if v.BinaryObject != nil {
		return string(v.BinaryObject.Value)
	}

	if v.BinaryGeometry != nil {
		return string(v.BinaryGeometry.Value)
	}

	if v.TextGeometry != nil {
		return v.TextGeometry.Value
	}

	if v.Decimal != nil {
		return v.Decimal.Value
	}

	if v.DateTime != nil {
		var year, month, day, hour, minute, second int32
		if v.DateTime.Year != nil {
			year = v.DateTime.Year.Int
		} else {
			year = 0
		}
		if v.DateTime.Month != nil {
			month = v.DateTime.Month.Int
		} else {
			month = 0
		}
		if v.DateTime.Day != nil {
			day = v.DateTime.Day.Int
		} else {
			day = 0
		}
		if v.DateTime.Hour != nil {
			hour = v.DateTime.Hour.Int
		} else {
			hour = 0
		}
		if v.DateTime.Minute != nil {
			minute = v.DateTime.Minute.Int
		} else {
			minute = 0
		}
		if v.DateTime.Second != nil {
			second = v.DateTime.Second.Int
		} else {
			second = 0
		}

		d := fmt.Sprintf("%.4d-%.2d-%.2d %.2d:%.2d:%.2d", year, month, day, hour, minute, second)
		parsed, _ := time.Parse("2006-01-02 15:04:05", d)

		return parsed.Format("2006-01-02 15:04:05")
	}

	return ""
}
