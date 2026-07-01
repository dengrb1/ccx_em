import { describe, expect, it } from 'vitest'
import { buildChannelPayload } from './channel-payload'

describe('buildChannelPayload', () => {
  it('Copilot 渠道省略 Base URL 时应写入默认上游地址', () => {
    const result = buildChannelPayload({
      name: 'copilot-channel',
      serviceType: 'copilot',
      baseUrl: '',
      baseUrls: [],
      website: '',
      insecureSkipVerify: false,
      lowQuality: false,
      injectDummyThoughtSignature: false,
      stripThoughtSignature: false,
      passbackReasoningContent: false,
      passbackThinkingBlocks: false,
      description: '',
      apiKeys: [],
      modelMapping: {},
      reasoningMapping: {},
      reasoningParamStyle: 'reasoning',
      textVerbosity: '',
      fastMode: false,
      customHeaders: {},
      proxyUrl: '',
      routePrefix: '',
      supportedModels: [],
      autoBlacklistBalance: true,
      normalizeMetadataUserId: true,
      stripEmptyTextBlocks: false,
      normalizeSystemRoleToTopLevel: false,
      codexNativeToolPassthrough: false,
      codexToolCompat: false,
      stripImageGenerationTool: false,
      noVision: false,
      noVisionModels: [],
      visionFallbackModel: ''
    })

    expect(result.baseUrl).toBe('https://api.githubcopilot.com')
    expect(result.baseUrls).toBeUndefined()
  })
})
