# AWS Trusted Advisor

1. [https://docs.aws.amazon.com/pdfs/trustedadvisor/latest/APIReference/ta-api.pdf](https://docs.aws.amazon.com/pdfs/trustedadvisor/latest/APIReference/ta-api.pdf)
1. [https://docs.aws.amazon.com/trustedadvisor/latest/APIReference/API_GetRecommendation.html](https://docs.aws.amazon.com/trustedadvisor/latest/APIReference/API_GetRecommendation.html)

## Troubleshooting

`2024/05/19 11:57:15 operation error TrustedAdvisor: ListRecommendations, https response error StatusCode: 403, RequestID: 11112222-3333-4444-5555-666677778888, AccessDeniedException: Access denied due to support level`

> How do I access AWS Trusted Advisor via API?
> 
> You can programmatically access Trusted Advisor best practice checks, recommendations, and prioritized recommendations using Trusted Advisor API, available to AWS Business Support, AWS Enterprise On-Ramp Support, and AWS Enterprise Support customers. To learn more about Trusted Advisor API, refer to the [user guide](https://docs.aws.amazon.com/awssupport/latest/user/get-started-with-aws-trusted-advisor-api.html).

Ref: [https://aws.amazon.com/premiumsupport/faqs/](https://aws.amazon.com/premiumsupport/faqs/)
 
> Note
> 
> * You must have a Business, Enterprise On-Ramp, or Enterprise Support plan to use the Trusted Advisor API
> * If you call the AWS Trusted Advisor API from an account that doesn't have a Business, Enterprise On-Ramp, or Enterprise Support plan, then you receive an Access Denied exception. For more information about changing your support plan, see AWS Support.

Ref: [https://docs.aws.amazon.com/awssupport/latest/user/get-started-with-aws-trusted-advisor-api.html](https://docs.aws.amazon.com/awssupport/latest/user/get-started-with-aws-trusted-advisor-api.html)