package command

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

type LDAPIDPWriteModel struct {
	eventstore.WriteModel

	ID                  string
	Name                string
	Host                string
	Port                string
	TLS                 bool
	BaseDN              string
	UserObjectClass     string
	UserUniqueAttribute string
	Admin               string
	Password            *crypto.CryptoValue
	idp.LDAPAttributes
	idp.Options

	State domain.IDPState
}

func (wm *LDAPIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Name = e.Name
			wm.Host = e.Host
			wm.Port = e.Port
			wm.TLS = e.TLS
			wm.BaseDN = e.BaseDN
			wm.UserObjectClass = e.UserObjectClass
			wm.UserUniqueAttribute = e.UserUniqueAttribute
			wm.Admin = e.Admin
			wm.Password = e.Password
			wm.LDAPAttributes = e.LDAPAttributes
			wm.Options = e.Options
			wm.State = domain.IDPStateActive
		case *idp.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *LDAPIDPWriteModel) reduceChangedEvent(e *idp.LDAPIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Host != nil {
		wm.Host = *e.Host
	}
	if e.Port != nil {
		wm.Port = *e.Port
	}
	if e.TLS != nil {
		wm.TLS = *e.TLS
	}
	if e.BaseDN != nil {
		wm.BaseDN = *e.BaseDN
	}
	if e.UserObjectClass != nil {
		wm.UserObjectClass = *e.UserObjectClass
	}
	if e.UserUniqueAttribute != nil {
		wm.UserUniqueAttribute = *e.UserUniqueAttribute
	}
	if e.Admin != nil {
		wm.Admin = *e.Admin
	}
	if e.Password != nil {
		wm.Password = e.Password
	}
	wm.LDAPAttributes.ReduceChanges(e.LDAPAttributeChanges)
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *LDAPIDPWriteModel) NewChanges(
	name,
	host,
	port string,
	tls bool,
	baseDN,
	userObjectClass,
	userUniqueAttribute,
	admin string,
	password string,
	secretCrypto crypto.Crypto,
	attributes idp.LDAPAttributes,
	options idp.Options,
) ([]idp.LDAPIDPChanges, error) {
	changes := make([]idp.LDAPIDPChanges, 0)
	var cryptedPassword *crypto.CryptoValue
	var err error
	if password != "" {
		cryptedPassword, err = crypto.Crypt([]byte(password), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeLDAPPassword(cryptedPassword))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeLDAPName(name))
	}
	if wm.Host != host {
		changes = append(changes, idp.ChangeLDAPHost(host))
	}
	if wm.Port != port {
		changes = append(changes, idp.ChangeLDAPPort(port))
	}
	if wm.TLS != tls {
		changes = append(changes, idp.ChangeLDAPTLS(tls))
	}
	if wm.BaseDN != baseDN {
		changes = append(changes, idp.ChangeLDAPBaseDN(baseDN))
	}
	if wm.UserObjectClass != userObjectClass {
		changes = append(changes, idp.ChangeLDAPUserObjectClass(userObjectClass))
	}
	if wm.UserUniqueAttribute != userUniqueAttribute {
		changes = append(changes, idp.ChangeLDAPUserUniqueAttribute(userUniqueAttribute))
	}
	if wm.Admin != admin {
		changes = append(changes, idp.ChangeLDAPAdmin(admin))
	}
	attrs := wm.LDAPAttributes.Changes(attributes)
	if !attrs.IsZero() {
		changes = append(changes, idp.ChangeLDAPAttributes(attrs))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeLDAPOptions(opts))
	}
	return changes, nil
}

type IDPRemoveWriteModel struct {
	eventstore.WriteModel

	ID    string
	State domain.IDPState
	name  string
}

func (wm *IDPRemoveWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.LDAPIDPAddedEvent:
			wm.reduceAdded(e.ID, e.Name)
		case *idp.LDAPIDPChangedEvent:
			wm.reduceChanged(e.ID, e.Name)
		case *idp.RemovedEvent:
			wm.reduceRemoved(e.ID)
		case *idpconfig.IDPConfigAddedEvent:
			wm.reduceAdded(e.ConfigID, "")
		case *idpconfig.IDPConfigRemovedEvent:
			wm.reduceRemoved(e.ConfigID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPRemoveWriteModel) reduceAdded(id string, name string) {
	if wm.ID != id {
		return
	}
	wm.State = domain.IDPStateActive
	wm.name = name
}

func (wm *IDPRemoveWriteModel) reduceChanged(id string, name *string) {
	if wm.ID != id || name == nil {
		return
	}
	wm.name = *name
}

func (wm *IDPRemoveWriteModel) reduceRemoved(id string) {
	if wm.ID != id {
		return
	}
	wm.State = domain.IDPStateRemoved
}
