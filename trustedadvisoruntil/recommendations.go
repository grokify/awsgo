package trustedadvisorutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/trustedadvisor"
)

func (svc TrustedAdvisorService) ListRecommendations(ctx context.Context, params *trustedadvisor.ListRecommendationsInput, optFns ...func(*trustedadvisor.Options)) (*trustedadvisor.ListRecommendationsOutput, error) {
	return svc.AWSSvcClient.ListRecommendations(ctx, params, optFns...)
}
