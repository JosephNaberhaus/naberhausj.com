function emulatorMain() {
    const fibProgram = '// A simple program that calculates the Fibonacci sequence.\n//\n// When the number exceeds the maximum value of an 8-bit number the program\n// will restart.\n//\n// "Register 1" will contain the previous Fibonacci number.\n// "Register 2" will contain the current Fibonacci number.\n\nstart:\nMV $i $1 1\nMV $i $2 1\nloop:\nADD $1 $2 $3\n// If the result is less than the previous value then we overflowed. Restart.\nLT $3 $2 $4\nJMPNZ $4 $i start\nMV $2 $1\nMV $3 $2\nJMP $i loop';
    const primeProgram = '// A program that calculates prime numbers.\n//\n// This program uses a simple, inefficient algorithm to check whether a number\n// is prime. For a number X it tries each number in the range [2, X - 1] and\n// checks whether X is evenly divided. If X is not evenly divided by any number\n// then X is a prime number.\n//\n// Since this algorithm takes progressively longer for larger numbers, you may\n// want to steadily increase the clock speed as it runs.\n//\n// "Register 4" will always contain the last discovered prime number.\n\n// Start with 2, the first prime number.\nMV $i $1 +2\nSTORE $i $1 0\n\nprime_loop:\n\n// Start dividing by 2.\nMV $i $3 +2\n\ncheck_divisable_loop:\nLOAD $i $2 0\n\n// Before we begin, check whether our divisor has reached our dividend.\n// If so, we\'ve found a prime number.\nEQ $2 $3 $1\nJMPNZ $1 $i is_prime\n\n// This loop uses repeated subtraction to perform division.\nsubtract_loop:\nSUB $2 $3 $2\n\n// If we\'re left with 0 then the number was evenly divided.\nEQ $2 $i $1 0\nJMPNZ $1 $i is_divisable\n\n// If the number is less than zero then the division resulted in a remainder.\nSLT $2 $i $1 0\nJMPNZ $1 $i not_divisable\n\n// Otherwise, we\'re not done dividing yet.\nJMP $i subtract_loop\n\n// The number was evenly divided which means the number is not prime.\n// Try the next number.\nis_divisable:\n// Also re-use this code if we simply want the next number.\nnext_number:\nLOAD $i $1 0\nADD $1 $i $1 1\nSTORE $i $1 0\n\nJMP $i prime_loop\n\n// The number wasn\'t evenly divided.\n// Try dividing by the next number.\nnot_divisable:\nADD $3 $i $3 1\nJMP $i check_divisable_loop\n\n// The number was prime. Display it in register 4.\nis_prime:\nLOAD $i $4 0\n\nJMP $i next_number';
    const sortProgram = '// A program that performs selection sort on a list of numbers.\n//\n// The numbers to sort are loaded in below. The program will sort these\n// numbers and then display the result in order.\n//\n// "Register 4" will show the numbers in sorted order when the algorithm is done.\n\n// Register 2 keeps track of how many numbers we have\nMV $i $2 0\n\n// Load some numbers into memory\nMV $i $1 10\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 4\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 100\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 23\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 1\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 42\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 45\nSTORE $2 $1\nADD $i $2 $2 1\n\nMV $i $1 245\nSTORE $2 $1\nADD $i $2 $2 1\n\n// Store register 2 in memory so that we can recall it\nSTORE $i $2 254\n\n// Memory 255 holds the current index\nMV $i $3 0\nSTORE $i $3 255\n\nsort_loop:\n\n// Register 3 holds the index of the smallest number so far\n// Register 4 holds the index of the current number\nLOAD $i $3 255\nMV $3 $4\n\nfind_smallest_loop:\n\nLOAD $3 $1\nLOAD $4 $2\n\nLT $1 $2 $1\nJMPNZ $1 $i not_smaller\n\n// The new number is smaller.\nMV $4 $3\n\nnot_smaller:\n\nADD $i $4 $4 1\n\n// Check if we\'re done\nLOAD $i $2 254\nEQ $4 $2 $1\nJMPNZ $1 $i done_find_smallest\n\nJMP $i find_smallest_loop\n\ndone_find_smallest:\n\n// Swap the numbers in memory\nLOAD $3 $1\nLOAD $i $4 255\nLOAD $4 $2\nSTORE $3 $2\nSTORE $4 $1\n\nADD $i $4 $4 1\nSTORE $i $4 255\n\n// Check if we\'re done sorting\nLOAD $i $2 254\nEQ $2 $4 $1\nJMPNZ $1 $i done_sorting\nJMP $i sort_loop\n\ndone_sorting:\nLOAD $i $2 254\n\n// The loop that displays the numbers in sorted order.\nstart_display:\nMV $i $1 0\n\ndisplay_loop:\nLOAD $1 $4\n\nADD $i $1 $1 1\nEQ $1 $2 $3\nJMPNZ $3 $i start_display\nJMP $i display_loop';

    const programByValue = {
        'fib': fibProgram,
        'prime': primeProgram,
        'sort': sortProgram
    };

    //
    // Preprocess
    //
    function getProProcessedProgram() {
        const programText = programInput.value;
        const lines = programText.split('\n');

        const labels = new Map();

        const program = [];
        for (let line of lines) {
            line = line.trim();
            if (line === '' || line.startsWith('//')) {
                continue;
            }

            const match = line.match(/^\s*([a-zA-Z_]+):\s*$/)
            if (match !== null && match.length === 2) {
                labels.set(match[1], program.length)
                continue
            }

            program.push(line);
        }

        const finalProgram = [];
        for (let line of program) {
            for (let entry of labels.entries()) {
                line = line.replaceAll(entry[0], `${entry[1]}`);
            }

            finalProgram.push(line);
        }

        return finalProgram;
    }

    //
    // Controls
    //
    let isRunning = false;
    let program = [];
    let rawProgramText = '';

    function start() {
        isRunning = true;

        let lastTime = document.timeline.currentTime;
        let fractionalFromLast = 0;
        function handleAnimationFrame(timeStamp) {
            if (!isRunning) {
                // Might as well render one more time.
                render();
                return;
            }

            const clockSpeedHz = 4 * clockSpeed.value;
            const numPulsesPerMillisecond = clockSpeedHz / 1000;
            const delta = timeStamp - lastTime;
            const numPulses = (delta * numPulsesPerMillisecond) + fractionalFromLast;
            const numPulsesToDo = Math.trunc(numPulses);
            fractionalFromLast = numPulses - numPulsesToDo;
            for (let i = 0; i < numPulsesToDo; i++) {
                pulse();
            }

            lastTime = timeStamp;
            requestAnimationFrame((t) => handleAnimationFrame(t));
        }

        requestAnimationFrame((t) => handleAnimationFrame(t));

        render();
    }

    function pause() {
        isRunning = false;

        render();
    }

    const stepButton = document.getElementById('step-button')
    stepButton.onclick = () => {
        pulse();
        render();
    }

    const resetButton = document.getElementById('reset-button')
    resetButton.onclick = () => {
        reset();
    }

    const startButton = document.getElementById('start-button')
    startButton.onclick = () => {
        start();
    }

    const pauseButton = document.getElementById('pause-button')
    pauseButton.onclick = () => {
        pause();
    }

    const clockSpeed = document.getElementById('clock-speed')
    clockSpeed.oninput = () => {
        let value = clockSpeed.value;

        if (value < 0) {
            value = 0;
        } else if (value > 5000) {
            value = 5000
        }

        clockSpeed.value = value;
    }

    const exampleProgramSelect = document.getElementById('load-program-select')

    const loadButton = document.getElementById('load-program-button');
    const onLoadButtonClick = () => {
        rawProgramText = programByValue[exampleProgramSelect.value];
        programInput.value = rawProgramText;
        program = getProProcessedProgram();
        reset();
    }
    loadButton.onclick = onLoadButtonClick;

    const programInput = document.getElementById('program-text');
    programInput.oninput = () => {
        rawProgramText = programInput.value;
        program = getProProcessedProgram();
        reset();
    }

    const errorMessage = document.getElementById('error-message');

    //
    // Emulator Code
    //
    const computerSVG = document.getElementById('computer');

    const registers = new Uint8Array(4);
    let registerBuffer = 0;

    const memory = new Uint8Array(256);

    const busInputs = new Uint8Array(2);
    let busOutput = 0;

    let programCounter = 0;
    let programCounterBuffer = 0;

    let clockStage = 0;

    function reset() {
        isRunning = false;

        function resetArray(arr) {
            for (let i = 0; i < arr.length; i++) {
                arr[i] = 0
            }
        }

        resetArray(registers)
        registerBuffer = 0;

        resetArray(memory);

        resetArray(busInputs);
        busOutput = 0;

        programCounter = 0;
        programCounterBuffer = 0

        errorMessage.style.setProperty('visibility', 'hidden');

        clockStage = -1;
        pulse();

        render();
    }

    function displayValueToLeds(value, displayID) {
        const displayGroup = computerSVG.getElementById(displayID)
        const leds = displayGroup.getElementsByTagName('circle');

        for (let i = 0; i < leds.length; i++) {
            const led = leds[i];
            const isOn = value & (0b1 << (leds.length - i - 1));

            if (isOn) {
                led.classList.add('on');
            } else {
                led.classList.remove('on');
            }
        }
    }

    function displayTextValue(value, displayID) {
        const textElement = computerSVG.getElementById(displayID)
        textElement.textContent = value;
    }

    let renderScheduled = false;
    function render() {
        displayValueToLeds(registers[0], 'register-1-display');
        displayValueToLeds(registers[1], 'register-2-display');
        displayValueToLeds(registers[2], 'register-3-display');
        displayValueToLeds(registers[3], 'register-4-display');
        displayValueToLeds(registerBuffer, 'register-buffer');
        displayValueToLeds(programCounter, 'program-counter');
        displayValueToLeds(programCounterBuffer, 'program-counter-buffer');
        displayValueToLeds(0b1 << clockStage, 'clock');
        displayValueToLeds(busInputs[0], 'bus-input-1');
        displayValueToLeds(busInputs[1], 'bus-input-2');
        displayValueToLeds(busOutput, 'bus-output');

        let line = currentLine();
        if (line.length > 18) {
            line = line.substring(0, 15) + '...';
        }
        displayTextValue(line, 'current-instruction')

        stepButton.disabled = isRunning;
        startButton.disabled = isRunning;
        pauseButton.disabled = !isRunning;
        programInput.disabled = isRunning;
        programInput.value = isRunning ? program.join('\n') : rawProgramText;

        renderScheduled = false;
    }

    function currentLine() {
        if (programCounter >= program.length) {
            return 'JMP $i 255';
        }

        return program[programCounter];
    }

    function currentLineParts() {
        return currentLine().split(/\s+/);
    }

    function assertNumArgumentsAllowImmediate(parts, numArgumentsWithoutImmediate) {
        if (parts.length === numArgumentsWithoutImmediate + 1) {
            return true
        }

        if (parts.length === numArgumentsWithoutImmediate + 2 && parts.includes('$i')) {
            return true
        }

        throw `"${parts[0]}" instructions should have ${numArgumentsWithoutImmediate} arguments (or ${numArgumentsWithoutImmediate + 1} when using an immediate value)`;
    }

    function to2sComplementForm(value) {
        if (value >= 0) {
            return value;
        }

        return (~(-value) + 1) & 0b11111111;
    }

    function from2sComplementForm(value) {
        if ((value & 0b10000000) === 0) {
            return value;
        }

        return -(~(value - 1) & 0b11111111);
    }

    function parseImmediate(immediateValue) {
        const is2s = immediateValue.startsWith('+') || immediateValue.startsWith('-');

        let parsedValue = parseInt(immediateValue, 10);
        if (isNaN(parsedValue)) {
            throw `invalid immediate value "${immediateValue}"`;
        }

        if (is2s) {
            if (parsedValue < -128 || parsedValue > 127) {
                throw 'signed immediate value is outside of the valid range [-128, 127]';
            }

            parsedValue = to2sComplementForm(parsedValue)
        } else  {
            if (parsedValue < 0 || parsedValue > 255) {
                throw 'unsigned immediate value is outside of the valid range [0, 255]';
            }
        }

        return parsedValue;
    }

    function evaluateBusInputs(parts, numInputs) {
        busInputs[0] = 0;
        busInputs[1] = 0;

        for (let i = 0; i < numInputs; i++) {
            switch (parts[i + 1]) {
                case '$1':
                    busInputs[i] = registers[0];
                    break;
                case '$2':
                    busInputs[i] = registers[1];
                    break;
                case '$3':
                    busInputs[i] = registers[2];
                    break;
                case '$4':
                    busInputs[i] = registers[3];
                    break;
                case '$i':
                    busInputs[i] = parseImmediate(parts[parts.length - 1]);
                    break;
                default:
                    throw `"${parts[i + 1]}" is not a valid input register selection`;
            }
        }
    }

    function evaluateInstruction() {
        const parts = currentLineParts();
        // Should never happen, but check it just in case.
        if (parts.length === 0) {
            throw 'line is empty';
        }

        switch (parts[0]) {
            case 'ADD':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] + busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'SUB':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] + to2sComplementForm(-busInputs[1])

                return {
                    outputRegister: parts[3],
                };
            case 'AND':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] & busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'OR':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] | busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'LT':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] < busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'SLT':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = from2sComplementForm(busInputs[0]) < from2sComplementForm(busInputs[1]);

                return {
                    outputRegister: parts[3],
                };
            case 'LTEQ':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] <= busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'SLTEQ':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = from2sComplementForm(busInputs[0]) <= from2sComplementForm(busInputs[1]);

                return {
                    outputRegister: parts[3],
                };
            case 'EQ':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] === busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'NEQ':
                assertNumArgumentsAllowImmediate(parts, 3);
                evaluateBusInputs(parts, 2);
                busOutput = busInputs[0] !== busInputs[1];

                return {
                    outputRegister: parts[3],
                };
            case 'MV':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 1);
                busOutput = busInputs[0];

                return {
                    outputRegister: parts[2],
                };
            case 'NOT':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 1);
                busOutput = (~busInputs[0]) & 0b11111111;

                return {
                    outputRegister: parts[2],
                };
            case 'SHIFTL':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 1);
                busOutput = (busInputs[0] << 1) & 0b11111111;

                return {
                    outputRegister: parts[2],
                };
            case 'SHIFTR':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 1);

                const hadLeading = busInputs[0] & 0b10000000;
                busOutput = (busInputs[0] >> 1) & 0b11111111;
                if (hadLeading) {
                    busOutput |= 0b10000000
                }

                return {
                    outputRegister: parts[2],
                };
            case 'JMP':
                assertNumArgumentsAllowImmediate(parts, 1);
                evaluateBusInputs(parts, 1);

                busOutput = 0;

                return {
                    isJump: true,
                };
            case 'JMPNZ':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 2);

                busOutput = 0;

                return {
                    isJumpNZ: true,
                };
            case 'LOAD':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 1);

                busOutput = memory[busInputs[0]]

                return {
                    outputRegister: parts[2],
                };
            case 'STORE':
                assertNumArgumentsAllowImmediate(parts, 2);
                evaluateBusInputs(parts, 2);

                busOutput = 0;

                return {
                    isMemoryWrite: true,
                };
            default:
                throw `"${parts[0]}" is not a valid operation`;
        }
    }

    function pulse() {
        clockStage++;
        if (clockStage === 4) {
            clockStage = 0;
        }

        try {
            switch (clockStage) {
                case 0: {
                    programCounter = programCounterBuffer;

                    const result = evaluateInstruction();

                    if (result.isMemoryWrite) {
                        memory[busInputs[0]] = busInputs[1];
                    }

                    registerBuffer = busOutput;

                    break;
                }
                case 1: {
                    const result = evaluateInstruction();

                    if (result.isJump) {
                        programCounterBuffer = busInputs[0];
                    } else if (result.isJumpNZ && busInputs[0] !== 0) {
                        programCounterBuffer = busInputs[1];
                    } else {
                        programCounterBuffer = programCounter + 1
                    }

                    break;
                }
                case 2: {
                    const result = evaluateInstruction();

                    if (result.outputRegister) {
                        switch (result.outputRegister) {
                            case '$1':
                                registers[0] = busOutput;
                                break;
                            case '$2':
                                registers[1] = busOutput;
                                break;
                            case '$3':
                                registers[2] = busOutput;
                                break;
                            case '$4':
                                registers[3] = busOutput;
                                break;
                            default:
                                throw `"${result.outputRegister}" is not a valid output register selection`;
                        }
                    }

                    // Need to evaluate again so that we render accurately with the new register value.
                    evaluateInstruction();

                    break;
                }
                case 3: {
                    evaluateInstruction();
                    break;
                }
            }
        } catch (message) {
            clockStage--;
            pause();
            displayErrorMessage(message);
        }

        render();
    }

    function displayErrorMessage(message) {
        errorMessage.textContent = message;
        errorMessage.style.setProperty('visibility', 'visible');
    }

    // Load an example program at startup.
    onLoadButtonClick();
}

emulatorMain()