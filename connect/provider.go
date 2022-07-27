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
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_CLIENT_CERTIFICATE_PATH", ""),
			},
			"ssl_client_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_CLIENT_KEY_PATH", ""),
			},
			"ssl_root_ca_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_ROOT_CA_PATH", ""),
			},
			"ssl_verify_server": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_VERIFY_SERVER", true),
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

	error := configureSSL(d, client)
	return client, error
}

func configureSSL(resources *schema.ResourceData, client kc.HighLevelClient) error {
	sslVerifyServer := resources.Get("ssl_verify_server").(bool)
	if !sslVerifyServer {
		// Must be done before any other ssl settings as this overrides the entire TLS Config
		client.SetInsecureSSL()
	}

	sslRootCaPath := resources.Get("ssl_root_ca_path").(string)
	if sslRootCaPath != "" {
		client.AddRootCertificate(sslRootCaPath)
	}

	clientCertificatePath := resources.Get("ssl_client_certificate_path").(string)
	clientCertificateKeyPath := resources.Get("ssl_client_key_path").(string)

	var error error = nil
	if clientCertificatePath != "" && clientCertificateKeyPath != "" {
		var certificate tls.Certificate
		certificate, error = tls.LoadX509KeyPair(clientCertificatePath, clientCertificateKeyPath)
		if error != nil {
			log.Printf("[ERROR] Loading certificate failed: %s", error)
		} else {
			client.SetClientCertificates(certificate)
		}
	}
	return error
}
