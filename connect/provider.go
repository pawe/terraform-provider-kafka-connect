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
			"ssl_client_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_CLIENT_CERTIFICATE", ""),
			},
			"ssl_client_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_CLIENT_KEY", ""),
			},
			"c": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_VERIFY_SERVER", true),
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

	address := d.Get("url").(string)
	client := kc.NewClient(address)

	user := d.Get("basic_auth_username").(string)
	pass := d.Get("basic_auth_password").(string)
	if user != "" && pass != "" {
		client.SetBasicAuth(user, pass)
	}

	verifyServer := d.Get("verify_server").(bool)
	if !verifyServer {
		client.SetInsecureSSL()
	}

	// client cert
	certificate := d.Get("ssl_client_certificate").(string)
	certificateKey := d.Get("ssl_client_key").(string)

	if certificate != "" && certificateKey != "" {
		certificate, error := tls.LoadX509KeyPair(certificate, certificateKey)
		if error != nil {
			log.Printf("[ERROR] Loading certificate failed: %s", error)
			return nil, error
		} else {
			client.SetClientCertificates(certificate)
		}
		return client, error
	}
	return client, nil
}
