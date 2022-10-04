package taskQueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
)

var defaultMessageLifeTime = 2

type TaskQueue struct {
	Name           string
	Path           string
	meta           meta
	client         *cloudtasks.Client
	ctx            context.Context
	initClientOnce sync.Once
	selfUrl        string
}

type meta interface {
	GetProjectId() string
	GetRegion() string
	GetEmail() string
}

func New(ctx context.Context, name string, meta meta, selfUrl string) *TaskQueue {
	return &TaskQueue{
		Name:           name,
		Path:           fmt.Sprintf("projects/%s/locations/%s/queues/%s", meta.GetProjectId(), meta.GetRegion(), name),
		meta:           meta,
		ctx:            ctx,
		initClientOnce: sync.Once{},
		selfUrl:        selfUrl,
	}
}

func (q *TaskQueue) EnqueueDeleteMessage(chatId int64, msgId int) error {

	if q.selfUrl == "" {
		log.Info("Self url is not set operation can't be enqued")
		return nil
	}

	var err error

	q.initClient()

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/chat/%d/%d", q.selfUrl, chatId, msgId)

	req := q.createDeleteRequest(url)

	createdTask, err := q.client.CreateTask(q.ctx, req)

	if err != nil {
		return err
	}

	log.Info("Message deletion task was created ", createdTask.Name)

	if err != nil {
		return fmt.Errorf("cloudtasks.CreateTask has failed: %v", err)
	}

	return nil
}

func (q *TaskQueue) Close() error {
	if q.client == nil {
		return nil
	}

	return q.client.Close()
}

func (q *TaskQueue) initClient() error {
	var err error

	init := func() error {
		client, err := cloudtasks.NewClient(q.ctx)
		if err != nil {
			return fmt.Errorf("cloudtasks client was not created: %v", err)
		}

		q.client = client

		return nil
	}

	q.initClientOnce.Do(func() {
		err = init()
	})

	return err
}

func (q *TaskQueue) createDeleteRequest(url string) *taskspb.CreateTaskRequest {
	req := &taskspb.CreateTaskRequest{
		Parent: q.Path,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_DELETE,
					Url:        url,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: q.meta.GetEmail(),
						},
					},
				},
			},
			ScheduleTime: getMinutesOffset(defaultMessageLifeTime),
		},
	}

	return req
}

func getMinutesOffset(minutes int) *timestamppb.Timestamp {
	d := time.Minute * time.Duration(minutes)

	ts := &timestamppb.Timestamp{
		Seconds: time.Now().Add(d).Unix(),
	}

	return ts
}
