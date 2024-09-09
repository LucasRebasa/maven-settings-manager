package structure

import "encoding/xml"

type Settings struct {
	XMLName           xml.Name  `xml:"settings"`
	Comment 				 string    `xml:",comment"`
	LocalRepository   string    `xml:"localRepository,omitempty"`
	InteractiveMode   string    `xml:"interactiveMode,omitempty"`
	Offline           string    `xml:"offline,omitempty"`
	PluginGroups      []string  `xml:"pluginGroups>pluginGroup,omitempty"`
	Servers           []Server  `xml:"servers>server,omitempty"`
	Mirrors           []Mirror  `xml:"mirrors>mirror,omitempty"`
	Proxies           []Proxy   `xml:"proxies>proxy,omitempty"`
	Profiles          []Profile `xml:"profiles>profile,omitempty"`
	ActiveProfiles    []string  `xml:"activeProfiles>activeProfile,omitempty"`
}

type Proxy struct {
	Id           string `xml:"id,omitempty"`
	Active       string `xml:"active,omitempty"`
	Protocol     string `xml:"protocol,omitempty"`
	Host         string `xml:"host,omitempty"`
	Port         string `xml:"port,omitempty"`
	Username     string `xml:"username,omitempty"`
	Password     string `xml:"password,omitempty"`
	NonProxyHost string `xml:"nonProxyHosts,omitempty"`
}

type Mirror struct {
	Id       string `xml:"id,omitempty"`
	Name     string `xml:"name,omitempty"`
	Url      string `xml:"url,omitempty"`
	MirrorOf string `xml:"mirrorOf,omitempty"`
}

type Server struct {
	Id                   string   `xml:"id,omitempty"`
	Username             string   `xml:"username,omitempty"`
	Password             string   `xml:"password,omitempty"`
	PrivateKey           string   `xml:"privateKey,omitempty"`
	Passphrase           string   `xml:"passphrase,omitempty"`
	FilePermissions      string   `xml:"filePermissions,omitempty"`
	DirectoryPermissions string   `xml:"directoryPermissions,omitempty"`
	Configuration        []AnyTag `xml:"configuration,omitempty"`
}

type Profile struct {
	XMLName            xml.Name     `xml:"profile,omitempty"`
	Id                 string       `xml:"id,omitempty"`
	Activation         Activation   `xml:"activation,omitempty"`
	Repositories       []Repository `xml:"repositories>repository,omitempty"`
	PluginRepositories []Repository `xml:"pluginRepositories>pluginRepository,omitempty"`
	Properties         AnyTag       `xml:"properties,omitempty"`
}

type AnyTag struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type Activation struct {
	ActiveByDefault string   `xml:"activeByDefault,omitempty"`
	Jdk             string   `xml:"jdk,omitempty"`
	Os              Os       `xml:"os,omitempty"`
	Property        Property `xml:"property,omitempty"`
	File            File     `xml:"file,omitempty"`
}

type File struct {
	Exists  string `xml:"exists,omitempty"`
	Missing string `xml:"missing,omitempty"`
}

type Property struct {
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

type Os struct {
	Name    string `xml:"name,omitempty"`
	Family  string `xml:"family,omitempty"`
	Arch    string `xml:"arch,omitempty"`
	Version string `xml:"version,omitempty"`
}

type Repository struct {
	Id        string      `xml:"id,omitempty"`
	Name      string      `xml:"name,omitempty"`
	Url       string      `xml:"url,omitempty"`
	Releases  ReleaseSnap `xml:"releases,omitempty"`
	Snapshots ReleaseSnap `xml:"snapshots,omitempty"`
	Layout    string      `xml:"layout,omitempty"`
}

type ReleaseSnap struct {
	Enabled        string `xml:"enabled,omitempty"`
	UpdatePolicy   string `xml:"updatePolicy,omitempty"`
	ChecksumPolicy string `xml:"checksumPolicy,omitempty"`
}
