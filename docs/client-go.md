# **client-go**

### 1、Client 客户端对象
（1）RESTClient客户端  
最基础的http请求，需要添加api。apis  
（2）DynamicClient客户端  
在RestClient上封装了一层，可以直接使用资源对象，不需要添加api。apis  
（3）ClientSet客户端  
在RestClient上封装了一层，可以直接使用资源对象，不需要添加api。apis  
（4）DiscoveryClient客户端  
获取k8s api server所支持的资源组，版本，资源对象等信息

---
### 2、Informer机制
#### （1）Reflector  
用于监控制定的k8s资源，当资源发生变化时，触发相应的变更事件  
#### （2）DeltaFIFO  
FIFO为一个先进先出的队列，他拥有队列操作的基本方法，Delta是一个资源对象存储，可以保存资源对象的操作类型  
#### （3）Indexer  
Indexer是client-go用户存储资源对象并自带索引功能的本地存储，例子IndexTest
* ThreadSafeMap并发安全存储
* Indexer索引器
* Indexer索引器核心实现 index.ByIndex   

---
### 3、WorkQueue
##### 主要的特性
* 有序
* 去重
* 并发性
* 标记机制
* 通知机制
* 延迟
* 限速
* Metric  
##### 三种队列
* Interface
* DelayingInterface
* RateLimitingInterface  
##### 三种队列的实现
* FIFO  
```go
//Queue
type Interface interface {
	Add(item interface{})
	Len() int
	Get() (item interface{}, shutdown bool)
	Done(item interface{})
	ShutDown()
	ShutDownWithDrain()
	ShuttingDown() bool
}
// FIFO队列数据结构
type Type struct {
    queue []t
    dirty set
    processing set
    cond *sync.Cond
    shuttingDown bool
    drain        bool
    metrics queueMetrics
    unfinishedWorkUpdatePeriod time.Duration
    clock                      clock.WithTicker
}

```
* DelayingQueue  
```go
type DelayingInterface interface {
	Interface
	AddAfter(item interface{}, duration time.Duration)
}
type delayingType struct {
    Interface
    clock clock.Clock
    stopCh chan struct{}
    stopOnce sync.Once
    heartbeat clock.Ticker
    waitingForAddCh chan *waitFor
    metrics retryMetrics
}
```
* RateLimitingQueue  
```go
type RateLimiter interface {
	When(item interface{}) time.Duration //获取指定元素应该等待的时间
	Forget(item interface{}) // 释放指定元素，清空该元素的排队数
	NumRequeues(item interface{}) int // 获取指定元素的排队数
}
```
##### 四种限速算法
* 排队指数算法（ItemExponentialFailureRateLimiter）
* 令牌桶算法（BucketRateLimiter）
* 计数器算法（ItemFastSlowRateLimiter）
* 混合模式算法（MaxOfRateLimiter）

---
### 4、EventBroadcaster事件管理器  
Kubernetes的事件（Event），是一种资源对象，用于展示集群内发生的情况，默认是一小时内
```go
type Event struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	InvolvedObject ObjectReference `json:"involvedObject" protobuf:"bytes,2,opt,name=involvedObject"`
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	Source EventSource `json:"source,omitempty" protobuf:"bytes,5,opt,name=source"`
	FirstTimestamp metav1.Time `json:"firstTimestamp,omitempty" protobuf:"bytes,6,opt,name=firstTimestamp"`
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty" protobuf:"bytes,7,opt,name=lastTimestamp"`
	Count int32 `json:"count,omitempty" protobuf:"varint,8,opt,name=count"`
	Type string `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`
	EventTime metav1.MicroTime `json:"eventTime,omitempty" protobuf:"bytes,10,opt,name=eventTime"`
	Series *EventSeries `json:"series,omitempty" protobuf:"bytes,11,opt,name=series"`
	Action string `json:"action,omitempty" protobuf:"bytes,12,opt,name=action"`
	Related *ObjectReference `json:"related,omitempty" protobuf:"bytes,13,opt,name=related"`
	ReportingController string `json:"reportingComponent" protobuf:"bytes,14,opt,name=reportingComponent"`
	ReportingInstance string `json:"reportingInstance" protobuf:"bytes,15,opt,name=reportingInstance"`
}
```
（1）EventRecorder：事件生产者，k8s系统组件通过EventRecorder记录关键性事件
```go
type EventRecorder interface {
	//对刚发生的事件进行记录
	Event(object runtime.Object, eventtype, reason, message string)
    //通过使用fmt.Sprintf格式化消息
	Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{})
    //功能和Eventf一样，但是可以添加annotations
	AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{})
}

```
（2）EventBroadcaster：事件消费者，消费EventRecorder记录的事件并将其分发给目前所有已链接的broadcasterWatcher  

两种分发机制
* 非阻塞分发机制 DropIfChannelFull
* 阻塞分发机制   WaitIfChannelFull
（3）broadcasterWatcher：观察者管理，用于定义事件的处理方式  
两种自定义处理事件的函数
* StartLogging： 将事件写入日志中
* StartRecordingToSink：将事件上报给Kubernetes API Server，并存储到ETCD集群


---
### 5、代码生成器
#### （1）client-gen  
用于生成ClientSet客户端的工具，通过Tags来识别是否需要生成
```makefile
make all WHAT=vendor/k8s.io/code-generator/cmd/client-gen
```
#### （2）lister-gen
一种为资源生成lister的工具，通过Tags来识别是否需要生成
```makefile
make all WHAT=vendor/k8s.io/code-generator/cmd/lister-gen
```
#### （3）informer-gen
一种为资源生成informer的工具，通过Tags来识别是否需要生成
```makefile
make all WHAT=vendor/k8s.io/code-generator/cmd/informer-gen
```


