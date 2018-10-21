package circleci

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jszwedko/go-circleci"
)

func resourceCircleciProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleciProjectCreate,
		Read:   resourceCircleciProjectRead,
		Update: resourceCircleciProjectUpdate,
		Delete: resourceCircleciProjectDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"env_vars": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"slack_chat_notifications_settings": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceCircleciProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client
	organization := meta.(*Organization).name
	project := d.Get("name").(string)
	envVars := d.Get("env_vars").(map[string]interface{})
	slackChatNofiticationsSettings := d.Get("slack_chat_notifications_settings").(map[string]interface{})

	_, err := client.FollowProject(organization, project)
	if err != nil {
		return err
	}

	d.SetId(project)

	for name, value := range envVars {
		_, err := client.AddEnvVar(organization, project, name, value.(string))
		if err != nil {
			return err
		}
	}

	if slackChatNofiticationsSettings != nil {
		nullOrString := func(key string) *string {
			if v, ok := slackChatNofiticationsSettings[key]; ok {
				r := v.(string)
				return &r
			}
			return nil
		}
		settings := circleci.SlackChatNotificationSettings{
			APIToken:        nullOrString("api_token"),
			Channel:         nullOrString("channel"),
			NotifyPrefs:     nullOrString("notify_prefs"),
			ChannelOverride: nullOrString("channel_override"),
			Subdomain:       nullOrString("subdomain"),
			WebhookURL:      nullOrString("webhook_url"),
		}
		if err := client.UpdateSlackChatNotificationsSettings(organization, project, settings); err != nil {
			return err
		}
	}

	return nil
}

func resourceCircleciProjectRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCircleciProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCircleciProjectCreate(d, meta)
}

func resourceCircleciProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client
	organization := meta.(*Organization).name
	name := d.Get("name").(string)

	return client.DisableProject(organization, name)
}
