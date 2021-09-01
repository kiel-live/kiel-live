package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCollectorStop(t *testing.T) {
	// given
	assert := assert.New(t)
	channelType := "stop"
	provider := "kvg"
	entityID := "123"
	channel := fmt.Sprintf("/%s/%s/%s", channelType, provider, entityID)

	// when
	collector, err := newCollector(nil, channel)

	// then
	if err != nil {
		t.Errorf("Channel type should not be supported: %v", err)
	}

	assert.Equal(channel, collector.channel, "should be equal")

	assert.Equal(channelType, collector.channelType, "should be equal")

	assert.Equal(provider, collector.provider, "should be equal")

	assert.Equal(entityID, collector.entityID, "should be equal")
}

func TestNewCollectorStops(t *testing.T) {
	// given
	assert := assert.New(t)
	channelType := "stops"
	channel := fmt.Sprintf("/%s", channelType)

	// when
	collector, err := newCollector(nil, channel)

	// then
	if err != nil {
		t.Errorf("Channel type should not be supported: %v", err)
	}

	assert.Equal(channel, collector.channel, "should be equal")

	assert.Equal(channelType, collector.channelType, "should be equal")
}

type MockedClient struct {
	mock.Mock
}

func (m *MockedClient) Publish(channel string, data string) error {
	args := m.Called(channel, data)
	return args.Error(0)
}

func TestRunStopsChannel(t *testing.T) {
	// given
	channel := "/stops"
	mockedClient := new(MockedClient)
	mockedClient.On("Publish", channel, mock.Anything).Return(nil)

	// when
	collector, _ := newCollector(mockedClient, channel)

	// then
	data := collector.collect()
	collector.publish(data)

	mockedClient.AssertExpectations(t)
}
