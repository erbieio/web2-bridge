package twitter

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweet/searchtweet"
	recenttypes "github.com/michimani/gotwi/tweet/searchtweet/types"
	"github.com/michimani/gotwi/tweet/timeline"
	tweettpyes "github.com/michimani/gotwi/tweet/timeline/types"
	"github.com/michimani/gotwi/user/userlookup"
	"github.com/michimani/gotwi/user/userlookup/types"
)

func GetUserinfoApi(c *gotwi.Client, text string) (resources.User, error) {
	p := &types.GetByUsernameInput{
		Username: text,
		UserFields: fields.UserFieldList{
			fields.UserFieldPublicMetrics,
			fields.UserFieldCreatedAt,
		},
	}
	res, err := userlookup.GetByUsername(context.Background(), c, p)
	if err != nil {
		return resources.User{}, err
	}

	return res.Data, nil
}

func GetUserTweetsListApi(c *gotwi.Client, id, pagination string) (*tweettpyes.ListTweetsOutput, error) {
	p := &tweettpyes.ListTweetsInput{
		ID:              id,
		MaxResults:      100,
		PaginationToken: pagination,
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldInReplyToUserID, fields.TweetFieldAuthorID, fields.TweetFieldText,
			fields.TweetFieldConversationID, fields.TweetFieldCreatedAt, fields.TweetFieldReferencedTweets,
			fields.TweetFieldPublicMetrics, fields.TweetFieldSource,
		},
	}
	res, err := timeline.ListTweets(context.Background(), c, p)
	if err != nil {
		return &tweettpyes.ListTweetsOutput{}, err
	}

	return res, nil
}

func GetUserinfoListApi(c *gotwi.Client, text []string) ([]resources.User, error) {
	p := &types.ListByUsernamesInput{
		Usernames: text,
		UserFields: fields.UserFieldList{
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}
	fmt.Println("request data=", *p)
	res, err := userlookup.ListByUsernames(context.Background(), c, p)
	if err != nil {
		return []resources.User{}, err
	}

	return res.Data, nil
}

type ClientOauth struct {
	oauthToken, oauthTokenSecret, apiKey, apiKeySecret string
}

func NewOAuth1Client(oauth *ClientOauth) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           oauth.oauthToken,
		OAuthTokenSecret:     oauth.oauthTokenSecret,
	}
	os.Setenv(gotwi.APIKeyEnvName, oauth.apiKey)
	os.Setenv(gotwi.APIKeySecretEnvName, oauth.apiKeySecret)
	c, err := gotwi.NewClient(in)
	if err != nil {
		return nil, err
	}

	c.SetOAuthConsumerKey(oauth.apiKey)
	signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(oauth.apiKeySecret), url.QueryEscape(oauth.oauthTokenSecret))
	c.SetSigningKey(signingKey)
	return c, nil
}

func NewOAuth2Client(oauth *ClientOauth) (*gotwi.Client, error) {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth2BearerToken,
	}
	os.Setenv(gotwi.APIKeyEnvName, oauth.apiKey)
	os.Setenv(gotwi.APIKeySecretEnvName, oauth.apiKeySecret)
	c, err := gotwi.NewClient(in)
	if err != nil {
		fmt.Println("NewOAuth2Client err =", err)
		return nil, err
	}
	c.SetOAuthConsumerKey(oauth.apiKey)
	signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(oauth.apiKeySecret), url.QueryEscape(oauth.oauthTokenSecret))
	c.SetSigningKey(signingKey)
	return c, err
}

func NewBearerTokenClient(accessToken string) (*gotwi.Client, error) {
	in := &gotwi.NewClientWithAccessTokenInput{
		AccessToken: accessToken,
	}
	c, err := gotwi.NewClientWithAccessToken(in)
	if err != nil {
		fmt.Println("NewOAuth2Client err =", err)
		return nil, err
	}
	return c, err
}

func SearchReplyRecentApi(c *gotwi.Client, text, pagination string) (*recenttypes.ListRecentOutput, error) {
	p := &recenttypes.ListRecentInput{
		//Query:      "keyword:" + text,
		Query:      `"` + text + `"`,
		MaxResults: 100,
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldInReplyToUserID, fields.TweetFieldAuthorID, fields.TweetFieldText,
			fields.TweetFieldConversationID, fields.TweetFieldCreatedAt, fields.TweetFieldReferencedTweets,
			fields.TweetFieldPublicMetrics, fields.TweetFieldSource,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldName, fields.UserFieldID, fields.UserFieldUrl, fields.UserFieldCreatedAt,
		},
		Expansions: fields.ExpansionList{
			fields.ExpansionAuthorID,
		},
	}
	if pagination != "" {
		p.NextToken = pagination
	}

	res, err := searchtweet.ListRecent(context.Background(), c, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SearchReplyRecentTweetReplyApi(c *gotwi.Client, text, pagination string) (*recenttypes.ListRecentOutput, error) {
	p := &recenttypes.ListRecentInput{
		Query:      "in_reply_to_tweet_id:" + text,
		MaxResults: 100,
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldInReplyToUserID, fields.TweetFieldAuthorID, fields.TweetFieldText,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldName, fields.UserFieldID, fields.UserFieldCreatedAt, fields.UserFieldPublicMetrics,
		},
		Expansions: fields.ExpansionList{
			fields.ExpansionAuthorID, fields.ExpansionInReplyToUserID,
		},
	}
	if pagination != "" {
		p.NextToken = pagination
	}

	res, err := searchtweet.ListRecent(context.Background(), c, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetTweetTimeLineMentions(c *gotwi.Client, id, pagination string) (*tweettpyes.ListMentionsOutput, error) {
	p := &tweettpyes.ListMentionsInput{
		ID:              id,
		MaxResults:      100,
		PaginationToken: pagination,
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldInReplyToUserID, fields.TweetFieldAuthorID, fields.TweetFieldText,
			fields.TweetFieldConversationID, fields.TweetFieldCreatedAt, fields.TweetFieldReferencedTweets,
			fields.TweetFieldPublicMetrics, fields.TweetFieldSource,
		},
	}

	res, err := timeline.ListMentions(context.Background(), c, p)
	if err != nil {
		return &tweettpyes.ListMentionsOutput{}, err
	}

	return res, nil
}
