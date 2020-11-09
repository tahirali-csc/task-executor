package scm

type Client struct{
	Git           GitService
	Webhooks      WebhookService
}