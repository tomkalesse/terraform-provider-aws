// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lexmodels

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelbuildingservice"
	awstypes "github.com/aws/aws-sdk-go-v2/service/lexmodelbuildingservice/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindBotVersionByName(ctx context.Context, conn *lexmodelbuildingservice.Client, name, version string) (*lexmodelbuildingservice.GetBotOutput, error) {
	input := &lexmodelbuildingservice.GetBotInput{
		Name:           aws.String(name),
		VersionOrAlias: aws.String(version),
	}

	output, err := conn.GetBot(ctx, input)

	if errs.IsA[*awstypes.NotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindSlotTypeVersionByName(ctx context.Context, conn *lexmodelbuildingservice.Client, name, version string) (*lexmodelbuildingservice.GetSlotTypeOutput, error) {
	input := &lexmodelbuildingservice.GetSlotTypeInput{
		Name:    aws.String(name),
		Version: aws.String(version),
	}

	output, err := conn.GetSlotType(ctx, input)

	if errs.IsA[*awstypes.NotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

// FindLatestBotVersionByName returns the latest published version of a bot or $LATEST if the bot has never been published.
// See https://docs.aws.amazon.com/lex/latest/dg/versioning-aliases.html.
func FindLatestBotVersionByName(ctx context.Context, conn *lexmodelbuildingservice.Client, name string) (string, error) {
	input := &lexmodelbuildingservice.GetBotVersionsInput{
		Name: aws.String(name),
	}
	var latestVersion int

	err := conn.GetBotVersionsPagesWithContext(ctx, input, func(page *lexmodelbuildingservice.GetBotVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, bot := range page.Bots {
			version := aws.ToString(bot.Version)

			if version == BotVersionLatest {
				continue
			}

			if version, err := strconv.Atoi(version); err != nil {
				continue
			} else if version > latestVersion {
				latestVersion = version
			}
		}

		return !lastPage
	})

	if err != nil {
		return "", err
	}

	if latestVersion == 0 {
		return BotVersionLatest, nil
	}

	return strconv.Itoa(latestVersion), nil
}

// FindLatestIntentVersionByName returns the latest published version of an intent or $LATEST if the intent has never been published.
// See https://docs.aws.amazon.com/lex/latest/dg/versioning-aliases.html.
func FindLatestIntentVersionByName(ctx context.Context, conn *lexmodelbuildingservice.Client, name string) (string, error) {
	input := &lexmodelbuildingservice.GetIntentVersionsInput{
		Name: aws.String(name),
	}
	var latestVersion int

	err := conn.GetIntentVersionsPagesWithContext(ctx, input, func(page *lexmodelbuildingservice.GetIntentVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, intent := range page.Intents {
			version := aws.ToString(intent.Version)

			if version == IntentVersionLatest {
				continue
			}

			if version, err := strconv.Atoi(version); err != nil {
				continue
			} else if version > latestVersion {
				latestVersion = version
			}
		}

		return !lastPage
	})

	if err != nil {
		return "", err
	}

	if latestVersion == 0 {
		return IntentVersionLatest, nil
	}

	return strconv.Itoa(latestVersion), nil
}

// FindLatestSlotTypeVersionByName returns the latest published version of a slot or $LATEST if the slot has never been published.
// See https://docs.aws.amazon.com/lex/latest/dg/versioning-aliases.html.
func FindLatestSlotTypeVersionByName(ctx context.Context, conn *lexmodelbuildingservice.Client, name string) (string, error) {
	input := &lexmodelbuildingservice.GetSlotTypeVersionsInput{
		Name: aws.String(name),
	}
	var latestVersion int

	err := conn.GetSlotTypeVersionsPagesWithContext(ctx, input, func(page *lexmodelbuildingservice.GetSlotTypeVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, slot := range page.SlotTypes {
			version := aws.ToString(slot.Version)

			if version == SlotTypeVersionLatest {
				continue
			}

			if version, err := strconv.Atoi(version); err != nil {
				continue
			} else if version > latestVersion {
				latestVersion = version
			}
		}

		return !lastPage
	})

	if err != nil {
		return "", err
	}

	if latestVersion == 0 {
		return SlotTypeVersionLatest, nil
	}

	return strconv.Itoa(latestVersion), nil
}
