package gosqs

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Admiral-Piett/goaws/app"
	"github.com/Admiral-Piett/goaws/app/models"
	"github.com/Admiral-Piett/goaws/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	utils.InitializeDecoders()
	m.Run()
}

func TestSendMessageBatch_POST_QueueNotFound(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "SendMessageBatch")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing")
	form.Add("SendMessageBatchRequestEntry.1.Id", "test_msg_001")
	form.Add("SendMessageBatchRequestEntry.1.MessageBody", "test%20message%20body%201")
	form.Add("SendMessageBatchRequestEntry.2.Id", "test_msg_002")
	form.Add("SendMessageBatchRequestEntry.2.MessageBody", "test%20message%20body%202")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageBatch)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := "NonExistentQueue"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSendMessageBatch_POST_NoEntry(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["testing"] = &app.Queue{Name: "testing"}

	form := url.Values{}
	form.Add("Action", "SendMessageBatch")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageBatch)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := "EmptyBatchRequest"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	req, _ = http.NewRequest("POST", "/", nil)
	form.Add("SendMessageBatchRequestEntry", "")
	req.PostForm = form

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSendMessageBatch_POST_IdNotDistinct(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["testing"] = &app.Queue{Name: "testing"}

	form := url.Values{}
	form.Add("Action", "SendMessageBatch")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing")
	form.Add("SendMessageBatchRequestEntry.1.Id", "test_msg_001")
	form.Add("SendMessageBatchRequestEntry.1.MessageBody", "test%20message%20body%201")
	form.Add("SendMessageBatchRequestEntry.2.Id", "test_msg_001")
	form.Add("SendMessageBatchRequestEntry.2.MessageBody", "test%20message%20body%202")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageBatch)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := "BatchEntryIdsNotDistinct"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSendMessageBatch_POST_TooManyEntries(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["testing"] = &app.Queue{Name: "testing"}

	form := url.Values{}
	form.Add("Action", "SendMessageBatch")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing")
	form.Add("SendMessageBatchRequestEntry.1.Id", "test_msg_001")
	form.Add("SendMessageBatchRequestEntry.1.MessageBody", "test%20message%20body%201")
	form.Add("SendMessageBatchRequestEntry.2.Id", "test_msg_002")
	form.Add("SendMessageBatchRequestEntry.2.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.3.Id", "test_msg_003")
	form.Add("SendMessageBatchRequestEntry.3.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.4.Id", "test_msg_004")
	form.Add("SendMessageBatchRequestEntry.4.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.5.Id", "test_msg_005")
	form.Add("SendMessageBatchRequestEntry.5.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.6.Id", "test_msg_006")
	form.Add("SendMessageBatchRequestEntry.6.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.7.Id", "test_msg_007")
	form.Add("SendMessageBatchRequestEntry.7.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.8.Id", "test_msg_008")
	form.Add("SendMessageBatchRequestEntry.8.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.9.Id", "test_msg_009")
	form.Add("SendMessageBatchRequestEntry.9.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.10.Id", "test_msg_010")
	form.Add("SendMessageBatchRequestEntry.10.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.11.Id", "test_msg_011")
	form.Add("SendMessageBatchRequestEntry.11.MessageBody", "test%20message%20body%202")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageBatch)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := "TooManyEntriesInBatchRequest"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSendMessageBatch_POST_Success(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["testing"] = &app.Queue{Name: "testing"}

	form := url.Values{}
	form.Add("Action", "SendMessageBatch")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing")
	form.Add("SendMessageBatchRequestEntry.1.Id", "test_msg_001")
	form.Add("SendMessageBatchRequestEntry.1.MessageBody", "test%20message%20body%201")
	form.Add("SendMessageBatchRequestEntry.2.Id", "test_msg_002")
	form.Add("SendMessageBatchRequestEntry.2.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.2.MessageAttribute.1.Name", "test_attribute_name_1")
	form.Add("SendMessageBatchRequestEntry.2.MessageAttribute.1.Value.StringValue", "test_attribute_value_1")
	form.Add("SendMessageBatchRequestEntry.2.MessageAttribute.1.Value.DataType", "String")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageBatch)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "<MD5OfMessageBody>1c538b76fce1a234bce865025c02b042</MD5OfMessageBody>"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSendMessageBatchToFIFOQueue_POST_Success(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["testing.fifo"] = &app.Queue{
		Name:   "testing.fifo",
		IsFIFO: true,
	}

	form := url.Values{}
	form.Add("Action", "SendMessageBatch")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing.fifo")
	form.Add("SendMessageBatchRequestEntry.1.Id", "test_msg_001")
	form.Add("SendMessageBatchRequestEntry.1.MessageGroupId", "GROUP-X")
	form.Add("SendMessageBatchRequestEntry.1.MessageBody", "test%20message%20body%201")
	form.Add("SendMessageBatchRequestEntry.2.Id", "test_msg_002")
	form.Add("SendMessageBatchRequestEntry.2.MessageGroupId", "GROUP-X")
	form.Add("SendMessageBatchRequestEntry.2.MessageBody", "test%20message%20body%202")
	form.Add("SendMessageBatchRequestEntry.2.MessageAttribute.1.Name", "test_attribute_name_1")
	form.Add("SendMessageBatchRequestEntry.2.MessageAttribute.1.Value.StringValue", "test_attribute_value_1")
	form.Add("SendMessageBatchRequestEntry.2.MessageAttribute.1.Value.DataType", "String")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageBatch)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "<MD5OfMessageBody>1c538b76fce1a234bce865025c02b042</MD5OfMessageBody>"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRequeueing_VisibilityTimeoutExpires(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//form := url.Values{}
	//form.Add("Action", "CreateQueue")
	//form.Add("QueueName", "requeue")
	//form.Add("Attribute.1.Name", "VisibilityTimeout")
	//form.Add("Attribute.1.Value", "1")
	//form.Add("Version", "2012-11-05")
	req.PostForm = url.Values{
		"Action":            []string{"CreateQueue"},
		"QueueName":         []string{"requeue"},
		"Attribute.1.Name":  []string{"VisibilityTimeout"},
		"Attribute.1.Value": []string{"1"},
		"Version":           []string{"2012-11-05"},
	}

	rr := httptest.NewRecorder()
	status, _ := CreateQueueV1(req)

	assert.Equal(t, http.StatusOK, status)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue")
	form.Add("MessageBody", "1")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()
	status, _ = SendMessageV1(req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()
	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	// try to receive another message.
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()
	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)
	time.Sleep(2 * time.Second)

	// message needs to be requeued
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	rr = httptest.NewRecorder()
	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)
	done <- struct{}{}
}

func TestRequeueing_ResetVisibilityTimeout(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "requeue-reset")
	form.Add("Attribute.1.Name", "VisibilityTimeout")
	form.Add("Attribute.1.Value", "10")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)
	assert.Equal(t, status, http.StatusOK)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset")
	form.Add("MessageBody", "1")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, resp := ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	receiptHandle := resp.GetResult().(models.ReceiveMessageResult).Messages[0].ReceiptHandle

	// try to receive another message.
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	// reset message visibility timeout to requeue it
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ChangeMessageVisibility")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset")
	form.Add("VisibilityTimeout", "0")
	form.Add("ReceiptHandle", receiptHandle)
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = ChangeMessageVisibilityV1(req)
	assert.Equal(t, status, http.StatusOK)

	// message needs to be requeued
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)
	done <- struct{}{}
}

func TestDeadLetterQueue(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	deadLetterQueue := &app.Queue{
		Name:     "failed-messages",
		Messages: []app.Message{},
	}
	app.SyncQueues.Lock()
	app.SyncQueues.Queues["failed-messages"] = deadLetterQueue
	app.SyncQueues.Unlock()
	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "testing-deadletter")
	form.Add("Attribute.1.Name", "VisibilityTimeout")
	form.Add("Attribute.1.Value", "1")
	form.Add("Attribute.2.Name", "RedrivePolicy")
	form.Add("Attribute.2.Value", `{"maxReceiveCount": 1, "deadLetterTargetArn":"arn:aws:sqs::000000000000:failed-messages"}`)
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)
	assert.Equal(t, status, http.StatusOK)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing-deadletter")
	form.Add("MessageBody", "1")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/testing-deadletter")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	time.Sleep(2 * time.Second)

	// receive the message one more time
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.PostForm = form

	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)
	time.Sleep(2 * time.Second)

	// another receive attempt
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.PostForm = form

	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)
	if len(deadLetterQueue.Messages) == 0 {
		t.Fatal("expected a message")
	}
}

func TestSetQueueAttributes_POST_QueueNotFound(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "SetQueueAttributes")
	form.Add("QueueUrl", "http://localhost:4100/queue/not-existing")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SetQueueAttributes)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := "NonExistentQueue"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSendingAndReceivingFromFIFOQueueReturnsSameMessageOnError(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "requeue-reset.fifo")
	form.Add("Attribute.1.Name", "VisibilityTimeout")
	form.Add("Attribute.1.Value", "2")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)
	assert.Equal(t, status, http.StatusOK)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset.fifo")
	form.Add("MessageBody", "1")
	form.Add("MessageGroupId", "GROUP-X")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset.fifo")
	form.Add("MessageBody", "2")
	form.Add("MessageGroupId", "GROUP-X")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset.fifo")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, resp := ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	result := resp.GetResult().(models.ReceiveMessageResult)
	receiptHandleFirst := result.Messages[0].ReceiptHandle
	if string(result.Messages[0].Body) != "1" {
		t.Fatalf("should have received body 1: %s", err)
	}

	// try to receive another message and we should get none
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset.fifo")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	if len(app.SyncQueues.Queues["requeue-reset.fifo"].FIFOMessages) != 1 {
		t.Fatal("there should be only 1 group locked")
	}

	if app.SyncQueues.Queues["requeue-reset.fifo"].FIFOMessages["GROUP-X"] != 0 {
		t.Fatal("there should be GROUP-X locked")
	}

	// remove message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "DeleteMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset.fifo")
	form.Add("ReceiptHandle", receiptHandleFirst)
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = DeleteMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	if len(app.SyncQueues.Queues["requeue-reset.fifo"].Messages) != 1 {
		t.Fatal("there should be only 1 message in queue")
	}

	// receive message - loop until visibility timeouts
	for {
		req, err = http.NewRequest("POST", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		form = url.Values{}
		form.Add("Action", "ReceiveMessage")
		form.Add("QueueUrl", "http://localhost:4100/queue/requeue-reset.fifo")
		form.Add("Version", "2012-11-05")
		req.PostForm = form

		status, resp := ReceiveMessageV1(req)
		assert.Equal(t, status, http.StatusOK)

		result := resp.GetResult().(models.ReceiveMessageResult)
		if len(result.Messages) == 0 {
			continue
		}

		if string(result.Messages[0].Body) != "2" {
			t.Fatalf("should have received body 2: %s", err)
		}
		break
	}

	done <- struct{}{}
}

func TestSendMessage_POST_DuplicatationNotAppliedToStandardQueue(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "stantdard-testing")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)

	assert.Equal(t, status, http.StatusOK)

	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/stantdard-testing")
	form.Add("MessageBody", "Test1")
	form.Add("MessageDeduplicationId", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)

	// Check the status code is what we expect.
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}
	if len(app.SyncQueues.Queues["stantdard-testing"].Messages) == 0 {
		t.Fatal("there should be 1 message in queue")
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/stantdard-testing")
	form.Add("MessageBody", "Test2")
	form.Add("MessageDeduplicationId", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)

	// Check the status code is what we expect.
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}
	if len(app.SyncQueues.Queues["stantdard-testing"].Messages) == 1 {
		t.Fatal("there should be 2 messages in queue")
	}
}

func TestSendMessage_POST_DuplicatationDisabledOnFifoQueue(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "no-dup-testing.fifo")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)

	assert.Equal(t, status, http.StatusOK)

	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/no-dup-testing.fifo")
	form.Add("MessageBody", "Test1")
	form.Add("MessageDeduplicationId", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)

	// Check the status code is what we expect.
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}
	if len(app.SyncQueues.Queues["no-dup-testing.fifo"].Messages) == 0 {
		t.Fatal("there should be 1 message in queue")
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/no-dup-testing.fifo")
	form.Add("MessageBody", "Test2")
	form.Add("MessageDeduplicationId", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)

	// Check the status code is what we expect.
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}
	if len(app.SyncQueues.Queues["no-dup-testing.fifo"].Messages) != 2 {
		t.Fatal("there should be 2 message in queue")
	}
}

func TestSendMessage_POST_DuplicatationEnabledOnFifoQueue(t *testing.T) {
	done := make(chan struct{}, 0)
	go PeriodicTasks(1*time.Second, done)

	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "dup-testing.fifo")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)

	assert.Equal(t, status, http.StatusOK)

	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.SyncQueues.Queues["dup-testing.fifo"].EnableDuplicates = true

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/dup-testing.fifo")
	form.Add("MessageBody", "Test1")
	form.Add("MessageDeduplicationId", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)

	// Check the status code is what we expect.
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}
	if len(app.SyncQueues.Queues["dup-testing.fifo"].Messages) == 0 {
		t.Fatal("there should be 1 message in queue")
	}

	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/dup-testing.fifo")
	form.Add("MessageBody", "Test2")
	form.Add("MessageDeduplicationId", "123")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)

	// Check the status code is what we expect.
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}
	if len(app.SyncQueues.Queues["dup-testing.fifo"].Messages) != 1 {
		t.Fatal("there should be 1 message in queue")
	}
	if body := app.SyncQueues.Queues["dup-testing.fifo"].Messages[0].MessageBody; string(body) == "Test2" {
		t.Fatal("duplicate message should not be added to queue")
	}
}

func TestSendMessage_POST_DelaySeconds(t *testing.T) {
	// create a queue
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{}
	form.Add("Action", "CreateQueue")
	form.Add("QueueName", "sendmessage-delay")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ := CreateQueueV1(req)

	assert.Equal(t, status, http.StatusOK)

	// send a message
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Add("Action", "SendMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/sendmessage-delay")
	form.Add("MessageBody", "1")
	form.Add("DelaySeconds", "2")
	form.Add("Version", "2012-11-05")
	req.PostForm = form

	status, _ = SendMessageV1(req)
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got \n%v want %v",
			status, http.StatusOK)
	}

	// receive message before delay is up
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/sendmessage-delay")
	form.Add("Version", "2012-11-05")
	req.PostForm = form
	status, _ = ReceiveMessageV1(req)
	assert.Equal(t, status, http.StatusOK)

	// receive message with wait should return after delay
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{}
	form.Add("Action", "ReceiveMessage")
	form.Add("QueueUrl", "http://localhost:4100/queue/sendmessage-delay")
	form.Add("WaitTimeSeconds", "10")
	form.Add("Version", "2012-11-05")
	req.PostForm = form
	start := time.Now()
	status, _ = ReceiveMessageV1(req)
	elapsed := time.Since(start)
	assert.Equal(t, status, http.StatusOK)
	if elapsed < 1*time.Second {
		t.Errorf("handler didn't wait at all")
	}
	if elapsed > 4*time.Second {
		t.Errorf("handler didn't need to wait all WaitTimeSeconds=10, only DelaySeconds=2")
	}
}

func TestCreateErrorResponseV1(t *testing.T) {
	expectedResponse := models.ErrorResponse{
		Result: models.ErrorResult{
			Type:    "Not Found",
			Code:    "AWS.SimpleQueueService.NonExistentQueue",
			Message: "The specified queue does not exist for this wsdl version.",
		},
		RequestId: "00000000-0000-0000-0000-000000000000",
	}
	status, response := createErrorResponseV1("QueueNotFound")

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, expectedResponse, response)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
// credits: https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
