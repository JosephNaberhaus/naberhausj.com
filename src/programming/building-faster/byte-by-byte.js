function byteByByteMain() {
    const inputOld = document.getElementById('byte-by-byte-input-old');
    const inputNew = document.getElementById('byte-by-byte-input-new');

    const displayOld = document.getElementById('byte-by-byte-display-old')
    const displayNew = document.getElementById('byte-by-byte-display-new')
    const result = document.getElementById('byte-by-byte-result');

    function onChange() {
        displayOld.innerHTML = '';
        displayNew.innerHTML = '';

        const oldValue = inputOld.value;
        const newValue = inputNew.value;

        for (let i = 0; i < oldValue.length; i++) {
            const letter = document.createElement('span');
            letter.innerText = '0x' + oldValue.charCodeAt(i).toString(16);
            letter.className = 'byte-by-byte-display-character ' + (i >= newValue.length ? 'not-match' : '');
            displayOld.appendChild(letter);
        }

        let allMatch = oldValue.length === newValue.length;
        for (let i = 0; i < newValue.length; i++) {
            let isMatch = true;
            if (i >= oldValue.length) {
                isMatch = false
            } else if (oldValue.charAt(i) !== newValue.charAt(i)) {
                isMatch = false;
            }

            if (!isMatch) {
                allMatch = false;
            }

            const letter = document.createElement('span');
            letter.innerText = '0x' + newValue.charCodeAt(i).toString(16);
            letter.className = 'byte-by-byte-display-character ' + (isMatch ? 'match' : 'not-match');
            displayNew.appendChild(letter);
        }

        result.value = allMatch ? 'Same' : 'Different'
    }

    inputOld.oninput = onChange;
    inputNew.oninput = onChange;

    onChange();
}

byteByByteMain();