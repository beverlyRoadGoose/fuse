package bot // import "heytobi.dev/fuse/bot"
import "testing"

func TestNew_ShouldBeAbleToInitializeBot(t *testing.T) {
	serviceProvider := &mockMessagingServiceProvider{}
	serviceProvider.On("SendMessage").Return(nil)
	_, _ = New(serviceProvider)
}
