function hashMain() {
    const inputOld = document.getElementById('hash-input-old');
    const inputNew = document.getElementById('hash-input-new');

    const displayOld = document.getElementById('hash-display-old')
    const displayNew = document.getElementById('hash-display-new')
    const result = document.getElementById('hash-result');

// Credit: https://stackoverflow.com/a/52171480/1768931
    const cyrb53 = function(str, seed = 0) {
        let h1 = 0xdeadbeef ^ seed, h2 = 0x41c6ce57 ^ seed;
        for (let i = 0, ch; i < str.length; i++) {
            ch = str.charCodeAt(i);
            h1 = Math.imul(h1 ^ ch, 2654435761);
            h2 = Math.imul(h2 ^ ch, 1597334677);
        }
        h1 = Math.imul(h1 ^ (h1>>>16), 2246822507) ^ Math.imul(h2 ^ (h2>>>13), 3266489909);
        h2 = Math.imul(h2 ^ (h2>>>16), 2246822507) ^ Math.imul(h1 ^ (h1>>>13), 3266489909);
        return 4294967296 * (2097151 & h2) + (h1>>>0);
    };

    const toStringFixedLength = function(number, radix, width) {
        let numberStr = number.toString(radix);
        while (numberStr.length < width) {
            numberStr = '0' + numberStr;
        }
        return numberStr;
    }

    function onChange() {
        displayOld.innerHTML = '';
        displayNew.innerHTML = '';

        const oldValue = inputOld.value;
        const newValue = inputNew.value;

        const oldValueHash = cyrb53(oldValue);
        const newValueHash = cyrb53(newValue);

        displayOld.innerText = '0x' + toStringFixedLength(oldValueHash, 16, 14);
        displayNew.innerText = '0x' + toStringFixedLength(newValueHash, 16, 14);
        displayNew.className = oldValueHash === newValueHash ? 'match' : 'not-match';

        result.value = oldValueHash === newValueHash ? 'Same' : 'Different'
    }

    inputOld.oninput = onChange;
    inputNew.oninput = onChange;

    onChange();
}

hashMain();