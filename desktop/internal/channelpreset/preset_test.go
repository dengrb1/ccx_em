package channelpreset

import "testing"

func TestBuildPayload(t *testing.T) {
	tests := []struct {
		name         string
		req          CreateChannelRequest
		wantTarget   string
		wantBaseURL  string
		wantService  string
		wantVision   bool
		wantPassback bool
	}{
		{
			name:        "deepseek messages",
			req:         CreateChannelRequest{Provider: ProviderDeepSeek, Target: TargetMessages, APIKey: "sk-test"},
			wantBaseURL: "https://api.deepseek.com/anthropic",
			wantService: "claude",
			wantVision:  true,
		},
		{
			name:         "mimo token plan",
			req:          CreateChannelRequest{Provider: ProviderMiMo, Target: TargetMessages, PlanID: "token-sgp", APIKey: "tp-test"},
			wantBaseURL:  "https://token-plan-sgp.xiaomimimo.com/v1",
			wantService:  "claude",
			wantPassback: true,
		},
		{
			name:        "kimi chat",
			req:         CreateChannelRequest{Provider: ProviderKimi, Target: TargetChat, APIKey: "sk-test"},
			wantBaseURL: "https://api.moonshot.cn/v1",
			wantService: "openai",
		},
		{
			name:        "glm chat",
			req:         CreateChannelRequest{Provider: ProviderGLM, Target: TargetChat, APIKey: "sk-test"},
			wantBaseURL: "https://open.bigmodel.cn/api/paas/v4",
			wantService: "openai",
		},
		{
			name:        "minimax chat",
			req:         CreateChannelRequest{Provider: ProviderMiniMax, Target: TargetChat, APIKey: "sk-test"},
			wantBaseURL: "https://api.minimax.chat/v1",
			wantService: "openai",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildPayload(tt.req)
			if err != nil {
				t.Fatalf("BuildPayload() error = %v", err)
			}
			if got.BaseURL != tt.wantBaseURL {
				t.Fatalf("BaseURL = %q, want %q", got.BaseURL, tt.wantBaseURL)
			}
			if got.ServiceType != tt.wantService {
				t.Fatalf("ServiceType = %q, want %q", got.ServiceType, tt.wantService)
			}
			if got.NoVision != tt.wantVision {
				t.Fatalf("NoVision = %v, want %v", got.NoVision, tt.wantVision)
			}
			if got.PassbackReasoningContent != tt.wantPassback {
				t.Fatalf("PassbackReasoningContent = %v, want %v", got.PassbackReasoningContent, tt.wantPassback)
			}
			if tt.req.Provider == ProviderMiMo {
				if got.ModelMapping["claude-sonnet-4-5"] != "mimo-v2.5-pro" {
					t.Fatalf("mimo model mapping missing: %#v", got.ModelMapping)
				}
				if got.VisionFallbackModel != "MiMo-V2.5" {
					t.Fatalf("VisionFallbackModel = %q, want MiMo-V2.5", got.VisionFallbackModel)
				}
			}
		})
	}
}

func TestBuildPayloadRejectsUnsupportedTarget(t *testing.T) {
	_, err := BuildPayload(CreateChannelRequest{Provider: ProviderKimi, Target: TargetMessages, APIKey: "sk-test"})
	if err == nil {
		t.Fatal("BuildPayload() expected error")
	}
}
