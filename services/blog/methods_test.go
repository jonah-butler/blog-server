package blog

import (
	"fmt"
	"os"
	"testing"
)

type ImageExtractionTest struct {
	Input string
	Want  int
}

func TestImageSourceExtraction(t *testing.T) {
	bucketName := os.Getenv("AWS_BUCKET")

	tests := []ImageExtractionTest{
		{
			Input: "<p>To get started with prettier in VSCode, first start off by installing the Prettier extension from the VSCode marketplace:</p><img src=\"https://dev-blog-resources.s3.amazonaws.com/Screenshot%202024-10-12%20at%2010.44.41%E2%80%AFPM.png\"><p>Then add the prettier development dependency to your project via:</p><pre class=\"ql-syntax\" data-language=\"shell\" spellcheck=\"false\"> npm install --save-dev --save-exact prettier</pre><p>Next add a <code class=\"inline-code\">.prettierrc</code> file to the root of your project.</p><p>This <code class=\"inline-code\">.prettierrc</code> file will contain all of the prettier configuration rules that will govern code formatting. A list of configuration can be found <a href=\"https://prettier.io/docs/en/options.html\">here on prettier's site</a>. Once you begin defining your configuration it will look similar to this simple example:</p><pre class=\"ql-syntax\" data-language=\"json\" spellcheck=\"false\"><span class=\"hljs-punctuation\">{</span>\r\n  <span class=\"hljs-attr\">&quot;singleQuote&quot;</span><span class=\"hljs-punctuation\">:</span> <span class=\"hljs-literal\"><span class=\"hljs-keyword\">true</span></span><span class=\"hljs-punctuation\">,</span>\r\n  <span class=\"hljs-attr\">&quot;semi&quot;</span><span class=\"hljs-punctuation\">:</span> <span class=\"hljs-literal\"><span class=\"hljs-keyword\">true</span></span>\r\n<span class=\"hljs-punctuation\">}</span></pre><p>At this point we will setup some configuration in VSCode to allow for auto-formatting a file whenever it is saved.</p><p>Go into the VSCode and locate the settings menu:</p><img src=\"https://dev-blog-resources.s3.amazonaws.com/Screenshot%202024-10-12%20at%2010.39.07%E2%80%AFPM.png\" data-caption=\"vscode on mac\"><p>Click into settings, and search the <code class=\"inline-code\">format</code> keyword within the <code class=\"inline-code\">User</code> tab:</p><img src=\"https://dev-blog-resources.s3.amazonaws.com/Screenshot%202024-10-12%20at%2010.39.38%E2%80%AFPM.png\" data-caption=\"vscode format settings\"><p>At this point we want to set Prettier as the default code formatter. In addition select the checkbox for <b>Editor: Format On Save.</b></p><p>Now that the general rules for functionality are set within our VSCode user, we need to enable some project specific configuration.</p><p>Next click the project tab to the right of the Workspace tab, in the example below, I have selected the dice-game project tab. Afterwards search the keyword <code class=\"inline-code\">prettier config</code>.</p><img src=\"https://dev-blog-resources.s3.amazonaws.com/Screenshot%202024-10-12%20at%2010.40.05%E2%80%AFPM.png\" data-caption=\"vscode prettier config settings\"><p>Within this pane, input the path to the .prettierrc file that was created earlier. If at the root of the project, simply put: <code class=\"inline-code\">./.prettierrc</code>. Next I like to require a prettier config in the option below the config path, so I keep this configuration on as well just to ensure no formatting will take place unless there is an explicit prettier file.</p><p>With all that under our belts, open a file, make some edits and save. It should be a voila moment, but if not, no troubles, sometimes a simple VSCode restart is required. That can be done easily with the command palette via:</p><p>Clicking: <code class=\"inline-code\">View/Command Palette</code></p><p>Or simply typing: <code class=\"inline-code\">CMD(Ctrl on Windows)+Shift+P</code><br></p><p>Type in <b>Reload Window </b>and select the option <b>Developer: Reload Window</b>. This will only take a quick second to reload, and your prettier configuration should now be complete!</p>",
			Want:  4,
		},
	}

	for _, test := range tests {
		imageSources, _ := extraImageSourcesFromHTML(test.Input, bucketName)
		fmt.Println(imageSources)

		if len(imageSources) != test.Want {
			t.Errorf("An invalid amount of images were extracted from the provided HTML input: wanted %d, got %d", test.Want, len(imageSources))
		}
	}

}

type ImageKeyExtractionTest struct {
	Input string
	Want  string
}

func TestExtractKeyFromImageSource(t *testing.T) {
	bucketName := os.Getenv("AWS_BUCKET")

	tests := []ImageKeyExtractionTest{
		{
			Input: "https://dev-blog-resources.s3.amazonaws.com/canvas_1736739686719.png",
			Want:  "canvas_1736739686719.png",
		},
	}

	for _, test := range tests {
		result := extractKeyFromImageSource(test.Input, bucketName+".s3.amazonaws.com/")

		if result != test.Want {
			t.Errorf("error extracting image key: wanted %s, got %s", test.Want, result)
		}
	}

}
