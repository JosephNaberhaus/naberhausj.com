function dependenciesMain() {
    const dependenciesGraph = document.getElementById('dependencies-graph');
    const results = document.getElementById('dependencies-results');

    const graph = new Map();
    graph.set('file-a', ['file-a', 'image-a', 'component-a', 'image-c']);
    graph.set('file-b', ['file-b', 'image-b', 'image-c']);
    graph.set('file-c', ['file-c', 'component-a', 'component-b', 'image-c']);
    graph.set('image-a', ['image-a']);
    graph.set('image-b', ['image-b']);
    graph.set('component-a', ['component-a', 'image-c']);
    graph.set('component-b', ['component-b']);
    graph.set('image-c', ['image-c']);

    const idToName = new Map();
    idToName.set('file-a', 'File A');
    idToName.set('file-b', 'File B');
    idToName.set('file-c', 'File C');
    idToName.set('image-a', 'Image A');
    idToName.set('image-b', 'Image B');
    idToName.set('component-a', 'Component A');
    idToName.set('component-b', 'Component B');
    idToName.set('image-c', 'Image C');

    const hasChanged = new Map();
    hasChanged.set('file-a', false);
    hasChanged.set('file-b', true);
    hasChanged.set('file-c', false);
    hasChanged.set('image-a', false);
    hasChanged.set('image-b', false);
    hasChanged.set('component-a', false);
    hasChanged.set('component-b', true);
    hasChanged.set('image-c', false);

    const rerender = function() {
        results.innerHTML = '';

        hasChanged.forEach((changed, id) => {
            dependenciesGraph.getElementById(id).setAttribute('fill', changed ? 'lightcoral' : 'lightgreen')

            for (let dependency of graph.get(id)) {
                if (hasChanged.get(dependency)) {
                    const result = document.createElement('div');
                    result.innerText = idToName.get(id);
                    results.appendChild(result);

                    break;
                }
            }
        });
    }

    hasChanged.forEach((_, id) => {
        dependenciesGraph.getElementById(id).onclick = () => {
            hasChanged.set(id, !hasChanged.get(id));
            rerender();
        };
    });

    rerender();
}

dependenciesMain();