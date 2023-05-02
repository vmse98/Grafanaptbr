package signature

import (
	"context"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
	openpgpErrors "github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/config"
	"github.com/grafana/grafana/pkg/plugins/manager/fakes"
	"github.com/grafana/grafana/pkg/plugins/manager/signature/statickey"
	"github.com/grafana/grafana/pkg/setting"
)

func TestReadPluginManifest(t *testing.T) {
	txt := `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

{
  "plugin": "grafana-googlesheets-datasource",
  "version": "1.0.0-dev",
  "files": {
    "LICENSE": "7df059597099bb7dcf25d2a9aedfaf4465f72d8d",
    "README.md": "08ec6d704b6115bef57710f6d7e866c050cb50ee",
    "gfx_sheets_darwin_amd64": "1b8ae92c6e80e502bb0bf2d0ae9d7223805993ab",
    "gfx_sheets_linux_amd64": "f39e0cc7344d3186b1052e6d356eecaf54d75b49",
    "gfx_sheets_windows_amd64.exe": "c8825dfec512c1c235244f7998ee95182f9968de",
    "module.js": "aaec6f51a995b7b843b843cd14041925274d960d",
    "module.js.LICENSE.txt": "7f822fe9341af8f82ad1b0c69aba957822a377cf",
    "module.js.map": "c5a524f5c4237f6ed6a016d43cd46938efeadb45",
    "plugin.json": "55556b845e91935cc48fae3aa67baf0f22694c3f"
  },
  "time": 1586817677115,
  "keyId": "7e4d0c6a708866e7"
}
-----BEGIN PGP SIGNATURE-----
Version: OpenPGP.js v4.10.1
Comment: https://openpgpjs.org

wqEEARMKAAYFAl6U6o0ACgkQfk0ManCIZuevWAIHSvcxOy1SvvL5gC+HpYyG
VbSsUvF2FsCoXUCTQflK6VdJfSPNzm8YdCdx7gNrBdly6HEs06ZaRp44F/ve
NR7DnB0CCQHO+4FlSPtXFTzNepoc+CytQyDAeOLMLmf2Tqhk2YShk+G/YlVX
74uuP5UXZxwK2YKJovdSknDIU7MhfuvvQIP/og==
=hBea
-----END PGP SIGNATURE-----`

	t.Run("valid manifest", func(t *testing.T) {
		s := ProvideService(&config.Cfg{}, statickey.New())
		manifest, err := s.readPluginManifest(context.Background(), []byte(txt))

		require.NoError(t, err)
		require.NotNil(t, manifest)
		assert.Equal(t, "grafana-googlesheets-datasource", manifest.Plugin)
		assert.Equal(t, "1.0.0-dev", manifest.Version)
		assert.Equal(t, int64(1586817677115), manifest.Time)
		assert.Equal(t, "7e4d0c6a708866e7", manifest.KeyID)
		expectedFiles := []string{"LICENSE", "README.md", "gfx_sheets_darwin_amd64", "gfx_sheets_linux_amd64",
			"gfx_sheets_windows_amd64.exe", "module.js", "module.js.LICENSE.txt", "module.js.map", "plugin.json",
		}
		assert.Equal(t, expectedFiles, fileList(manifest))
	})

	t.Run("invalid manifest", func(t *testing.T) {
		modified := strings.ReplaceAll(txt, "README.md", "xxxxxxxxxx")
		s := ProvideService(&config.Cfg{}, statickey.New())
		_, err := s.readPluginManifest(context.Background(), []byte(modified))
		require.Error(t, err)
	})
}

func TestReadPluginManifestV2(t *testing.T) {
	txt := `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

{
  "manifestVersion": "2.0.0",
  "signatureType": "private",
  "signedByOrg": "willbrowne",
  "signedByOrgName": "Will Browne",
  "rootUrls": [
    "http://localhost:3000/"
  ],
  "plugin": "test",
  "version": "1.0.0",
  "time": 1605807018050,
  "keyId": "7e4d0c6a708866e7",
  "files": {
    "plugin.json": "2bb467c0bfd6c454551419efe475b8bf8573734e73c7bab52b14842adb62886f"
  }
}
-----BEGIN PGP SIGNATURE-----
Version: OpenPGP.js v4.10.1
Comment: https://openpgpjs.org

wqIEARMKAAYFAl+2q6oACgkQfk0ManCIZudmzwIJAXWz58cd/91rTXszKPnE
xbVEvERCbjKTtPBQBNQyqEvV+Ig3MuBSNOVy2SOGrMsdbS6lONgvgt4Cm+iS
wV+vYifkAgkBJtg/9DMB7/iX5O0h49CtSltcpfBFXlGqIeOwRac/yENzRzAA
khdr/tZ1PDgRxMqB/u+Vtbpl0xSxgblnrDOYMSI=
=rLIE
-----END PGP SIGNATURE-----`

	t.Run("valid manifest", func(t *testing.T) {
		s := ProvideService(&config.Cfg{}, statickey.New())
		manifest, err := s.readPluginManifest(context.Background(), []byte(txt))

		require.NoError(t, err)
		require.NotNil(t, manifest)
		assert.Equal(t, "test", manifest.Plugin)
		assert.Equal(t, "1.0.0", manifest.Version)
		assert.Equal(t, int64(1605807018050), manifest.Time)
		assert.Equal(t, "7e4d0c6a708866e7", manifest.KeyID)
		assert.Equal(t, "2.0.0", manifest.ManifestVersion)
		assert.Equal(t, plugins.PrivateSignature, manifest.SignatureType)
		assert.Equal(t, "willbrowne", manifest.SignedByOrg)
		assert.Equal(t, "Will Browne", manifest.SignedByOrgName)
		assert.Equal(t, []string{"http://localhost:3000/"}, manifest.RootURLs)
		assert.Equal(t, []string{"plugin.json"}, fileList(manifest))
	})
}

func TestCalculate(t *testing.T) {
	t.Run("Validate root URL against App URL for non-private plugin if is specified in manifest", func(t *testing.T) {
		tcs := []struct {
			appURL            string
			expectedSignature plugins.Signature
		}{
			{
				appURL: "https://dev.grafana.com",
				expectedSignature: plugins.Signature{
					Status:     plugins.SignatureValid,
					Type:       plugins.GrafanaSignature,
					SigningOrg: "Grafana Labs",
				},
			},
			{
				appURL: "https://non.matching.url.com",
				expectedSignature: plugins.Signature{
					Status: plugins.SignatureInvalid,
				},
			},
		}

		parentDir, err := filepath.Abs("../")
		if err != nil {
			t.Errorf("could not construct absolute path of current dir")
			return
		}

		for _, tc := range tcs {
			origAppURL := setting.AppUrl
			t.Cleanup(func() {
				setting.AppUrl = origAppURL
			})
			setting.AppUrl = tc.appURL

			basePath := filepath.Join(parentDir, "testdata/non-pvt-with-root-url/plugin")
			s := ProvideService(&config.Cfg{}, statickey.New())
			sig, err := s.Calculate(context.Background(), &fakes.FakePluginSource{
				PluginClassFunc: func(ctx context.Context) plugins.Class {
					return plugins.External
				},
			}, plugins.FoundPlugin{
				JSONData: plugins.JSONData{
					ID: "test-datasource",
					Info: plugins.Info{
						Version: "1.0.0",
					},
				},
				FS: mustNewStaticFSForTests(t, basePath),
			})
			require.NoError(t, err)
			require.Equal(t, tc.expectedSignature, sig)
		}
	})

	t.Run("Unsigned Chromium file should not invalidate signature for Renderer plugin running on Windows", func(t *testing.T) {
		backup := runningWindows
		t.Cleanup(func() {
			runningWindows = backup
		})

		basePath := "../testdata/renderer-added-file/plugin"

		runningWindows = true
		s := ProvideService(&config.Cfg{}, statickey.New())
		sig, err := s.Calculate(context.Background(), &fakes.FakePluginSource{
			PluginClassFunc: func(ctx context.Context) plugins.Class {
				return plugins.External
			},
		}, plugins.FoundPlugin{
			JSONData: plugins.JSONData{
				ID:   "test-renderer",
				Type: plugins.Renderer,
				Info: plugins.Info{
					Version: "1.0.0",
				},
			},
			FS: mustNewStaticFSForTests(t, basePath),
		})
		require.NoError(t, err)
		require.Equal(t, plugins.Signature{
			Status:     plugins.SignatureValid,
			Type:       plugins.GrafanaSignature,
			SigningOrg: "Grafana Labs",
		}, sig)
	})

	t.Run("Signature verification should work with any path separator", func(t *testing.T) {
		const basePath = "../testdata/app-with-child/dist"

		platformWindows := fsPlatform{separator: '\\'}
		platformUnix := fsPlatform{separator: '/'}

		type testCase struct {
			name      string
			platform  fsPlatform
			fsFactory func() (plugins.FS, error)
		}
		var testCases []testCase
		for _, fsFactory := range []struct {
			name string
			f    func() (plugins.FS, error)
		}{
			{"local fs", func() (plugins.FS, error) {
				return plugins.NewLocalFS(basePath), nil
			}},
			{"static fs", func() (plugins.FS, error) {
				return plugins.NewStaticFS(plugins.NewLocalFS(basePath))
			}},
		} {
			testCases = append(testCases, []testCase{
				{"unix " + fsFactory.name, platformUnix, fsFactory.f},
				{"windows " + fsFactory.name, platformWindows, fsFactory.f},
			}...)
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Replace toSlash for cross-platform testing
				oldToSlash := toSlash
				oldFromSlash := fromSlash
				t.Cleanup(func() {
					toSlash = oldToSlash
					fromSlash = oldFromSlash
				})
				toSlash = tc.platform.toSlashFunc()
				fromSlash = tc.platform.fromSlashFunc()

				s := ProvideService(&config.Cfg{}, statickey.New())
				pfs, err := tc.fsFactory()
				require.NoError(t, err)
				pfs, err = newPathSeparatorOverrideFS(string(tc.platform.separator), pfs)
				require.NoError(t, err)
				sig, err := s.Calculate(context.Background(), &fakes.FakePluginSource{
					PluginClassFunc: func(ctx context.Context) plugins.Class {
						return plugins.External
					},
				}, plugins.FoundPlugin{
					JSONData: plugins.JSONData{
						ID:   "myorgid-simple-app",
						Type: plugins.App,
						Info: plugins.Info{
							Version: "%VERSION%",
						},
					},
					FS: pfs,
				})
				require.NoError(t, err)
				require.Equal(t, plugins.Signature{
					Status:     plugins.SignatureValid,
					Type:       plugins.GrafanaSignature,
					SigningOrg: "Grafana Labs",
				}, sig)
			})
		}
	})
}

type fsPlatform struct {
	separator rune
}

// toSlashFunc returns a new function that acts as filepath.ToSlash but for the specified os-separator.
// This can be used to test filepath.ToSlash-dependant code cross-platform.
func (p fsPlatform) toSlashFunc() func(string) string {
	return func(path string) string {
		if p.separator == '/' {
			return path
		}
		return strings.ReplaceAll(path, string(p.separator), "/")
	}
}

// fromSlashFunc returns a new function that acts as filepath.FromSlash but for the specified os-separator.
// This can be used to test filepath.FromSlash-dependant code cross-platform.
func (p fsPlatform) fromSlashFunc() func(string) string {
	return func(path string) string {
		if p.separator == '/' {
			return path
		}
		return strings.ReplaceAll(path, "/", string(p.separator))
	}
}

func TestFsPlatform(t *testing.T) {
	t.Run("unix", func(t *testing.T) {
		toSlashUnix := fsPlatform{'/'}.toSlashFunc()
		require.Equal(t, "folder", toSlashUnix("folder"))
		require.Equal(t, "/folder", toSlashUnix("/folder"))
		require.Equal(t, "/folder/file", toSlashUnix("/folder/file"))
		require.Equal(t, "/folder/other\\file", toSlashUnix("/folder/other\\file"))
	})

	t.Run("windows", func(t *testing.T) {
		toSlashWindows := fsPlatform{'\\'}.toSlashFunc()
		require.Equal(t, "folder", toSlashWindows("folder"))
		require.Equal(t, "C:/folder", toSlashWindows("C:\\folder"))
		require.Equal(t, "folder/file.exe", toSlashWindows("folder\\file.exe"))
	})
}

// fsPathSeparatorFiles embeds a plugins.FS and overrides the Files() behaviour so all the returned elements
// have the specified path separator. This can be used to test Files() behaviour cross-platform.
type fsPathSeparatorFiles struct {
	plugins.FS

	separator string
}

// newPathSeparatorOverrideFS returns a new fsPathSeparatorFiles. Sep is the separator that will be used ONLY for
// the elements returned by Files().
func newPathSeparatorOverrideFS(sep string, ufs plugins.FS) (fsPathSeparatorFiles, error) {
	return fsPathSeparatorFiles{
		FS:        ufs,
		separator: sep,
	}, nil
}

// Files returns LocalFS.Files(), but all path separators for the current platform (filepath.Separator)
// are replaced with f.separator.
func (f fsPathSeparatorFiles) Files() ([]string, error) {
	files, err := f.FS.Files()
	if err != nil {
		return nil, err
	}
	const osSepStr = string(filepath.Separator)
	for i := 0; i < len(files); i++ {
		files[i] = strings.ReplaceAll(files[i], osSepStr, f.separator)
	}
	return files, nil
}

func (f fsPathSeparatorFiles) Open(name string) (fs.File, error) {
	return f.FS.Open(strings.ReplaceAll(name, f.separator, string(filepath.Separator)))
}

func TestFSPathSeparatorFiles(t *testing.T) {
	for _, tc := range []struct {
		name string
		sep  string
	}{
		{"unix", "/"},
		{"windows", "\\"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pfs, err := newPathSeparatorOverrideFS(
				"/", plugins.NewInMemoryFS(
					map[string][]byte{"a": nil, strings.Join([]string{"a", "b", "c"}, tc.sep): nil},
				),
			)
			require.NoError(t, err)
			files, err := pfs.Files()
			require.NoError(t, err)
			exp := []string{"a", strings.Join([]string{"a", "b", "c"}, tc.sep)}
			sort.Strings(files)
			sort.Strings(exp)
			require.Equal(t, exp, files)
		})
	}
}

func fileList(manifest *PluginManifest) []string {
	var keys []string
	for k := range manifest.Files {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Test_urlMatch_privateGlob(t *testing.T) {
	type args struct {
		specs  []string
		target string
	}
	tests := []struct {
		name        string
		args        args
		shouldMatch bool
	}{
		{
			name: "Support single wildcard matching single subdomain",
			args: args{
				specs:  []string{"https://*.example.com"},
				target: "https://test.example.com",
			},
			shouldMatch: true,
		},
		{
			name: "Do not support single wildcard matching multiple subdomains",
			args: args{
				specs:  []string{"https://*.example.com"},
				target: "https://more.test.example.com",
			},
			shouldMatch: false,
		},
		{
			name: "Support multiple wildcards matching multiple subdomains",
			args: args{
				specs:  []string{"https://**.example.com"},
				target: "https://test.example.com",
			},
			shouldMatch: true,
		},
		{
			name: "Support multiple wildcards matching multiple subdomains",
			args: args{
				specs:  []string{"https://**.example.com"},
				target: "https://more.test.example.com",
			},
			shouldMatch: true,
		},
		{
			name: "Support single wildcard matching single paths",
			args: args{
				specs:  []string{"https://www.example.com/*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: true,
		},
		{
			name: "Do not support single wildcard matching multiple paths",
			args: args{
				specs:  []string{"https://www.example.com/*"},
				target: "https://www.example.com/other/grafana",
			},
			shouldMatch: false,
		},
		{
			name: "Support double wildcard matching multiple paths",
			args: args{
				specs:  []string{"https://www.example.com/**"},
				target: "https://www.example.com/other/grafana",
			},
			shouldMatch: true,
		},
		{
			name: "Do not support subdomain mismatch",
			args: args{
				specs:  []string{"https://www.test.example.com/grafana/docs"},
				target: "https://www.dev.example.com/grafana/docs",
			},
			shouldMatch: false,
		},
		{
			name: "Support single wildcard matching single path",
			args: args{
				specs:  []string{"https://www.example.com/grafana*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: true,
		},
		{
			name: "Do not support single wildcard matching different path prefix",
			args: args{
				specs:  []string{"https://www.example.com/grafana*"},
				target: "https://www.example.com/somethingelse",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support path mismatch",
			args: args{
				specs:  []string{"https://example.com/grafana"},
				target: "https://example.com/grafana1",
			},
			shouldMatch: false,
		},
		{
			name: "Support both domain and path wildcards",
			args: args{
				specs:  []string{"https://*.example.com/*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: true,
		},
		{
			name: "Do not support wildcards without TLDs",
			args: args{
				specs:  []string{"https://example.*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: false,
		},
		{
			name: "Support exact match",
			args: args{
				specs:  []string{"https://example.com/test"},
				target: "https://example.com/test",
			},
			shouldMatch: true,
		},
		{
			name: "Does not support scheme mismatch",
			args: args{
				specs:  []string{"https://test.example.com/grafana"},
				target: "http://test.example.com/grafana",
			},
			shouldMatch: false,
		},
		{
			name: "Support trailing slash in spec",
			args: args{
				specs:  []string{"https://example.com/"},
				target: "https://example.com",
			},
			shouldMatch: true,
		},
		{
			name: "Support trailing slash in target",
			args: args{
				specs:  []string{"https://example.com"},
				target: "https://example.com/",
			},
			shouldMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := urlMatch(tt.args.specs, tt.args.target, plugins.PrivateGlobSignature)
			require.NoError(t, err)
			require.Equal(t, tt.shouldMatch, got)
		})
	}
}

func Test_urlMatch_private(t *testing.T) {
	type args struct {
		specs  []string
		target string
	}
	tests := []struct {
		name        string
		args        args
		shouldMatch bool
	}{
		{
			name: "Support exact match",
			args: args{
				specs:  []string{"https://example.com/test"},
				target: "https://example.com/test",
			},
			shouldMatch: true,
		},
		{
			name: "Support trailing slash in spec",
			args: args{
				specs:  []string{"https://example.com/test/"},
				target: "https://example.com/test",
			},
			shouldMatch: true,
		},
		{
			name: "Support trailing slash in target",
			args: args{
				specs:  []string{"https://example.com/test"},
				target: "https://example.com/test/",
			},
			shouldMatch: true,
		},
		{
			name: "Do not support single wildcard matching single subdomain",
			args: args{
				specs:  []string{"https://*.example.com"},
				target: "https://test.example.com",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support multiple wildcards matching multiple subdomains",
			args: args{
				specs:  []string{"https://**.example.com"},
				target: "https://more.test.example.com",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support single wildcard matching single paths",
			args: args{
				specs:  []string{"https://www.example.com/*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support double wildcard matching multiple paths",
			args: args{
				specs:  []string{"https://www.example.com/**"},
				target: "https://www.example.com/other/grafana",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support subdomain mismatch",
			args: args{
				specs:  []string{"https://www.test.example.com/grafana/docs"},
				target: "https://www.dev.example.com/grafana/docs",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support path mismatch",
			args: args{
				specs:  []string{"https://example.com/grafana"},
				target: "https://example.com/grafana1",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support both domain and path wildcards",
			args: args{
				specs:  []string{"https://*.example.com/*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support wildcards without TLDs",
			args: args{
				specs:  []string{"https://example.*"},
				target: "https://www.example.com/grafana1",
			},
			shouldMatch: false,
		},
		{
			name: "Do not support scheme mismatch",
			args: args{
				specs:  []string{"https://test.example.com/grafana"},
				target: "http://test.example.com/grafana",
			},
			shouldMatch: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := urlMatch(tt.args.specs, tt.args.target, plugins.PrivateSignature)
			require.NoError(t, err)
			require.Equal(t, tt.shouldMatch, got)
		})
	}
}

func Test_validateManifest(t *testing.T) {
	tcs := []struct {
		name        string
		manifest    *PluginManifest
		expectedErr string
	}{
		{
			name:        "Empty plugin field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.Plugin = "" }),
			expectedErr: "valid manifest field plugin is required",
		},
		{
			name:        "Empty keyId field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.KeyID = "" }),
			expectedErr: "valid manifest field keyId is required",
		},
		{
			name:        "Empty signedByOrg field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.SignedByOrg = "" }),
			expectedErr: "valid manifest field signedByOrg is required",
		},
		{
			name:        "Empty signedByOrgName field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.SignedByOrgName = "" }),
			expectedErr: "valid manifest field SignedByOrgName is required",
		},
		{
			name:        "Empty signatureType field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.SignatureType = "" }),
			expectedErr: "valid manifest field signatureType is required",
		},
		{
			name:        "Invalid signatureType field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.SignatureType = "invalidSignatureType" }),
			expectedErr: "valid manifest field signatureType is required",
		},
		{
			name:        "Empty files field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.Files = map[string]string{} }),
			expectedErr: "valid manifest field files is required",
		},
		{
			name:        "Empty time field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.Time = 0 }),
			expectedErr: "valid manifest field time is required",
		},
		{
			name:        "Empty version field",
			manifest:    createV2Manifest(t, func(m *PluginManifest) { m.Version = "" }),
			expectedErr: "valid manifest field version is required",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s := ProvideService(&config.Cfg{}, statickey.New())
			err := s.validateManifest(context.Background(), *tc.manifest, nil)
			require.Errorf(t, err, tc.expectedErr)
		})
	}
}

func createV2Manifest(t *testing.T, cbs ...func(*PluginManifest)) *PluginManifest {
	t.Helper()

	m := &PluginManifest{
		Plugin:  "grafana-test-app",
		Version: "2.5.3",
		KeyID:   "7e4d0c6a708866e7",
		Time:    1586817677115,
		Files: map[string]string{
			"plugin.json": "55556b845e91935cc48fae3aa67baf0f22694c3f",
		},
		ManifestVersion: "2.0.0",
		SignatureType:   plugins.GrafanaSignature,
		SignedByOrg:     "grafana",
		SignedByOrgName: "grafana",
	}

	for _, cb := range cbs {
		cb(m)
	}

	return m
}

func mustNewStaticFSForTests(t *testing.T, dir string) plugins.FS {
	sfs, err := plugins.NewStaticFS(plugins.NewLocalFS(dir))
	require.NoError(t, err)
	return sfs
}

type revokedKeyProvider struct{}

func (p *revokedKeyProvider) GetPublicKey(ctx context.Context, keyID string) (string, error) {
	// dummy revoked key created locally
	const publicKeyText = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
YNE=
=U6y9
-----END PGP PUBLIC KEY BLOCK-----
`
	return publicKeyText, nil
}

func Test_VerifyRevokedKey(t *testing.T) {
	s := ProvideService(&config.Cfg{}, &revokedKeyProvider{})
	m := createV2Manifest(t)
	txt := `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

{
  "manifestVersion": "2.0.0",
  "signatureType": "grafana",
  "signedByOrg": "grafana",
  "signedByOrgName": "Grafana Labs",
  "plugin": "test-app",
  "version": "1.0.0",
  "time": 1621356785895,
  "keyId": "7e4d0c6a708866e7",
  "files": {
	"plugin.json": "c59a51bf6d7ecd7a99608ccb99353390c8b973672a938a0247164324005c0caf",
	"dashboards/connections.json": "bea86da4be970b98dc4681802ab55cdef3441dc3eb3c654cb207948d17b25303",
	"dashboards/memory.json": "7c042464941084caa91d0a9a2f188b05315a9796308a652ccdee31ca4fbcbfee",
	"dashboards/connections_result.json": "124d85c9c2e40214b83273f764574937a79909cfac3f925276fbb72543c224dc"
  }
}
-----BEGIN PGP SIGNATURE-----

iQGzBAEBCgAdFiEE7JeUrw5MX53LZ8CoXfB349i2H4YFAmRKigIACgkQXfB349i2
H4ZdKgwAuVuTjGT7Rn1MfxYRUXRymdnyqsDRYaK8gw5i9OZweBuJBVLtL1eFII0h
tTr+2jM4kGlsCakpJm3sjRG//8sBYoO5GsnOM6g1gv7mgUwo/Pv3A5eFFeOIkF1W
E33nNyF17BlY+YPVJPMQ8Q4uBSz2pDlcdQY8gOleWERWMWvmsHZgobt7wyGgts7Y
hCzKdm+e5/HpWBskW7dRMh1yB+8Ql+IK/Ksy8EDdX+Yv1fGV6ZNNIQxSEBXSily6
uvZlU9zExa0db9rkg53jFpSfSFpQIJJ0Y0yOmHKDA4WLnphroCIBwo2lxIBIwuNH
sXjmTjacvrqk13Af7Gat7XSNLapBfy5rTZwJFOwGWyDP1V0FTrlmt5vmoD0MRskq
gry5NAKktwc2llGaS5uGc5wJ1wTvl5wYQkU8lBevdejntpQSOYNEuICe+OyKQP+h
OOKpUCovEat+3W9JU1PM+z3cb1H/WWQ3hpKEykyzzi/jZMuRnRobW8Jm/4WxFgaY
70RA9/V8
=NUH5
-----END PGP SIGNATURE-----
`
	block, _ := clearsign.Decode([]byte(txt))
	require.NotNil(t, block, "failed to decode block")
	err := s.validateManifest(context.Background(), *m, block)
	require.Error(t, err)
	require.Contains(t, err.Error(), openpgpErrors.ErrKeyRevoked.Error())
}
