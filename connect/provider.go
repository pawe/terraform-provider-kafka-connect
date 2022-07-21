package connect

import (
	"crypto/tls"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kc "github.com/ricardo-ch/go-kafka-connect/lib/connectors"
)

func Provider() *schema.Provider {
	log.Printf("[INFO] Creating Provider")
	provider := schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_URL", ""),
			},
			"basic_auth_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_BASIC_AUTH_USERNAME", ""),
			},
			"basic_auth_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_BASIC_AUTH_PASSWORD", ""),
			},
			"ssl_client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_CLIENT_CERTIFICATE", ""),
			},
			"ssl_client_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_CLIENT_KEY", ""),
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"kafka-connect_connector": kafkaConnectorResource(),
		},
	}
	log.Printf("[INFO] Created provider: %v", provider)
	return &provider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing KafkaConnect client")
	addr := d.Get("url").(string)
	c := kc.NewClient(addr)
	user := d.Get("basic_auth_username").(string)
	pass := d.Get("basic_auth_password").(string)
	if user != "" && pass != "" {
		c.SetBasicAuth(user, pass)
	}

	// certs
	cert := d.Get("ssl_client_certificate_path").(string)
	key := d.Get("ssl_client_key_path").(string)

	if len(cert) > 0 && len(key) > 0 {
		cert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			log.Fatalf("client: loadkeys: %s", err)
		} else {
			c.SetClientCertificates(cert)
		}
	}

	return c, nil
}
