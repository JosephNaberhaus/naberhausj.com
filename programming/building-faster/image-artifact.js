function imageArtifactMain(){const e=document.getElementById('image-artifact-input'),n=document.getElementById('image-artifact-image'),s=document.getElementById('image-artifact-output'),o=[100,250,500,750,1e3,1500,2e3];function t(){const t=+e.value;n.width=t;const i=[];for(let e of o){if(e>t){i.push({Width:t,Height:t,File:`image-${t}.png`});break}if(i.push({Width:e,Height:e,File:`image-${e}.png`}),t===e)break}const a={OriginalWidth:t,OriginalHeight:t,Files:i};s.innerText=JSON.stringify(a,null,'	')}e.oninput=t,t()}imageArtifactMain()