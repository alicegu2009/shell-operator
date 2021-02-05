package hook

import (
	"github.com/flant/shell-operator/pkg/app"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/flant/shell-operator/pkg/hook/controller"
	"github.com/flant/shell-operator/pkg/webhook/conversion"
	"github.com/flant/shell-operator/pkg/webhook/validating"
	. "github.com/flant/shell-operator/pkg/webhook/validating/types"
)

func newHookManager(t *testing.T, testdataDir string) (*hookManager, func()) {
	var err error
	hm := NewHookManager()
	tmpDir, err := ioutil.TempDir("", "hook_manager")
	if err != nil {
		t.Fatalf("Make tmpdir should not fail: %v", err)
	}
	hooksDir, _ := filepath.Abs(testdataDir)
	hm.WithDirectories(hooksDir, tmpDir)
	conversionManager := conversion.NewWebhookManager()
	conversionManager.Settings = app.ConversionWebhookSettings
	hm.WithConversionWebhookManager(conversionManager)
	validatingManager := validating.NewWebhookManager()
	validatingManager.Settings = app.ValidatingWebhookSettings
	hm.WithValidatingWebhookManager(validatingManager)

	return hm, func() { os.RemoveAll(tmpDir) }
}

func Test_HookManager_Init(t *testing.T) {
	hooksDir := "testdata/hook_manager"
	hm, rmFn := newHookManager(t, hooksDir)
	defer rmFn()

	if !strings.HasSuffix(hm.WorkingDir(), hooksDir) {
		t.Fatalf("Hook manager should has working dir '%s', got: '%s'", hooksDir, hm.WorkingDir())
	}

	err := hm.Init()
	if err != nil {
		t.Fatalf("Hook manager Init should not fail: %v", err)
	}
}

func Test_HookManager_GetHookNames(t *testing.T) {
	hm, rmFn := newHookManager(t, "testdata/hook_manager")
	defer rmFn()

	err := hm.Init()
	if err != nil {
		t.Fatalf("Hook manager Init should not fail: %v", err)
	}

	names := hm.GetHookNames()

	expectedCount := 4
	if len(names) != expectedCount {
		t.Fatalf("Hook manager should have %d hooks, got %d", expectedCount, len(names))
	}

	// TODO fix sorting!!!
	expectedNames := []string{
		"configMapHooks/hook.sh",
		"hook.sh",
		"podHooks/hook.sh",
		"podHooks/hook2.sh",
	}

	for i, expectedName := range expectedNames {
		if names[i] != expectedName {
			t.Fatalf("Hook manager should have hook '%s' at index %d, %s", expectedName, i, names[i])
		}
	}

}

func TestHookController_HandleValidatingEvent(t *testing.T) {
	g := NewWithT(t)

	hm, rmFn := newHookManager(t, "testdata/hook_manager_validating")
	defer rmFn()

	err := hm.Init()
	if err != nil {
		t.Fatalf("Hook manager Init should not fail: %v", err)
	}

	ev := ValidatingEvent{
		WebhookId:       "ololo-policy-example-com",
		ConfigurationId: "hooks",
		Review:          nil,
	}

	h := hm.GetHook("hook.sh")
	h.HookController.EnableValidatingBindings()

	canHandle := h.HookController.CanHandleValidatingEvent(ev)

	g.Expect(canHandle).To(BeTrue())

	var infoList []controller.BindingExecutionInfo
	h.HookController.HandleValidatingEvent(ev, func(info controller.BindingExecutionInfo) {
		infoList = append(infoList, info)
	})

	g.Expect(infoList).Should(HaveLen(1))

}

func Test_HookManager_conversion_chains(t *testing.T) {
	g := NewWithT(t)

	hm, rmFn := newHookManager(t, "testdata/hook_manager_conversion_chains")
	defer rmFn()

	err := hm.Init()
	g.Expect(err).ShouldNot(HaveOccurred(), "Hook manager Init should not fail: %v", err)

	// conversion chain for 1 crd
	g.Expect(hm.conversionChains).Should(HaveLen(2))

	crdName := "crontabs.stable.example.com"
	g.Expect(hm.conversionChains).Should(HaveKey(crdName))

	chain := hm.conversionChains[crdName]
	// 6 paths in cache for each binding.
	g.Expect(chain.PathsCache).Should(HaveLen(6))

	var convPath []string

	// Find path for unknown crd
	convPath = hm.FindConversionChain("unknown"+crdName, conversion.ConversionRule{
		FromVersion: "azaza",
		ToVersion:   "ololo",
	})
	g.Expect(convPath).Should(BeNil())

	// Find path for unknown from version
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "unknown-version",
		ToVersion:   "ololo",
	})
	g.Expect(convPath).Should(BeNil())

	// Find path for unknown to version
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "ololo",
		ToVersion:   "unknown-version",
	})
	g.Expect(convPath).Should(BeNil())

	// Find path for unknown from and to versions
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "from-unknown-version",
		ToVersion:   "to-unknown-version",
	})
	g.Expect(convPath).Should(BeNil())

	// Find a simple path.
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "azaza",
		ToVersion:   "ololo",
	})
	g.Expect(convPath).Should(HaveLen(1))

	// Find a full path in an "up" direction.
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "azaza",
		ToVersion:   "abc",
	})
	g.Expect(convPath).Should(HaveLen(3))
	g.Expect(convPath[0]).Should(Equal("azaza->ololo"))
	g.Expect(convPath[1]).Should(Equal("ololo->foobar"))
	g.Expect(convPath[2]).Should(Equal("foobar->abc"))

	// Find a full path in a "down" direction.
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "abc",
		ToVersion:   "azaza",
	})
	g.Expect(convPath).Should(HaveLen(3))
	g.Expect(convPath[0]).Should(Equal("abc->foobar"))
	g.Expect(convPath[1]).Should(Equal("foobar->ololo"))
	g.Expect(convPath[2]).Should(Equal("ololo->azaza"))

	// Find a part path in an "up" direction.
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "ololo",
		ToVersion:   "abc",
	})
	g.Expect(convPath).Should(HaveLen(2))
	g.Expect(convPath[0]).Should(Equal("ololo->foobar"))
	g.Expect(convPath[1]).Should(Equal("foobar->abc"))

	// Find a part path in a "down" direction.
	convPath = hm.FindConversionChain(crdName, conversion.ConversionRule{
		FromVersion: "foobar",
		ToVersion:   "azaza",
	})
	g.Expect(convPath).Should(HaveLen(2))
	g.Expect(convPath[0]).Should(Equal("foobar->ololo"))
	g.Expect(convPath[1]).Should(Equal("ololo->azaza"))

	// Cache has 6 paths from bindings, 2 more paths for each full path and 1 more path for each part path.
	g.Expect(chain.PathsCache).Should(HaveLen(6 + 2 + 2 + 1 + 1))
}
