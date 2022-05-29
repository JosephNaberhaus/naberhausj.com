function imageArtifactMain() {
    const input = document.getElementById('image-artifact-input');
    const image = document.getElementById('image-artifact-image');
    const output = document.getElementById('image-artifact-output');

    const outputWidth = [100, 250, 500, 750, 1000, 1500, 2000];

    function onChange() {
        const width = +input.value;

        image.width = width;

        const files = [];
        for (let targetWidth of outputWidth) {
            if (targetWidth > width) {
                files.push({
                    Width: width,
                    Height: width,
                    File: `image-${width}.png`,
                });
                break;
            }

            files.push({
                Width: targetWidth,
                Height: targetWidth,
                File: `image-${targetWidth}.png`,
            });

            if (width === targetWidth) {
                break;
            }
        }

        const artifact = {
            OriginalWidth: width,
            OriginalHeight: width,
            Files: files,
        }

        output.innerText = JSON.stringify(artifact, null, '\t');
    }

    input.oninput = onChange;

    onChange();
}

imageArtifactMain();