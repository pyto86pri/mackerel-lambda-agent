# [WIP] mackerel-agent-lambda

A monitoring agent of [Mackerel](https://mackerel.io/) for AWS Lambda.

This is [AWS Lambda Extensions](https://aws.amazon.com/jp/blogs/compute/introducing-aws-lambda-extensions-in-preview/) provided as AWS Lambda Layers.

**NOTE: This is experimental and not suitable for production use.**

- Deploy layer with AWS SAM CLI.
```console
$ sam build
$ sam deploy --guided
```
- Set up the layer for Lambda functions which you want to monitor on Mackerel.
- Configure Mackerel API key as `MACKEREL_API_KEY` in the Lambda functions environment variables.

## Note
- AWS Lambda Extensions runs on the same execution environment as Lambda functions. So it can impact function performance.
- Following are overhead estimations;
  - CPU overhead          : ?
  - Memory overhead       : ~70MB (Working on decreasing)
  - Duration time overhead: ~1msec