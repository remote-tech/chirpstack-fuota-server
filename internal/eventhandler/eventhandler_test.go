package eventhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"

	"github.com/brocaar/lorawan/applayer/clocksync"
	"github.com/remote-tech/chirpstack-fuota-server/v4/internal/client/as"
	"github.com/remote-tech/chirpstack-fuota-server/v4/internal/test"
	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/chirpstack/chirpstack/api/go/v4/integration"
)

func TestEventHandler(t *testing.T) {
	assert := require.New(t)

	handler, err := NewHandler(HandlerOptions{
		JSON: false,
	})
	assert.NoError(err)

	s := httptest.NewServer(handler)
	defer s.Close()

	t.Run("ClockSync", func(t *testing.T) {
		assert := require.New(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := test.NewMockDeviceServiceClient(ctrl)
		as.SetDeviceClient(mock)

		cmdAns := clocksync.Command{
			CID: clocksync.AppTimeAns,
			Payload: &clocksync.AppTimeAnsPayload{
				TimeCorrection: 10,
				Param: clocksync.AppTimeAnsPayloadParam{
					TokenAns: 123,
				},
			},
		}
		cmdAnsB, err := cmdAns.MarshalBinary()
		assert.NoError(err)

		mock.EXPECT().Enqueue(
			gomock.Any(),
			gomock.Eq(&api.EnqueueDeviceQueueItemRequest{
				QueueItem: &api.DeviceQueueItem{
					DevEui: "0102030405060708",
					FPort:  202,
					Data:   cmdAnsB,
				}}),
		).Return(nil, nil)

		cmd := clocksync.Command{
			CID: clocksync.AppTimeReq,
			Payload: &clocksync.AppTimeReqPayload{
				DeviceTime: 200,
				Param: clocksync.AppTimeReqPayloadParam{
					AnsRequired: true,
					TokenReq:    123,
				},
			},
		}
		cmdB, err := cmd.MarshalBinary()
		assert.NoError(err)

		uplink := integration.UplinkEvent{
			DeviceInfo: &integration.DeviceInfo{
				DevEui: "0102030405060708",
			},
			Data:  cmdB,
			FPort: uint32(clocksync.DefaultFPort),
			RxInfo: []*gw.UplinkRxInfo{
				{
					TimeSinceGpsEpoch: ptypes.DurationProto(time.Second * 210),
				},
			},
		}
		uplinkB, err := proto.Marshal(&uplink)
		assert.NoError(err)

		resp, err := http.Post(s.URL+"?event=up", "application/protobuf", bytes.NewBuffer(uplinkB))
		assert.NoError(err)
		assert.Equal(200, resp.StatusCode)
		time.Sleep(time.Millisecond * 50) // request is handled async
	})
}
