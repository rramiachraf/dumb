const fullAbout = document.querySelector("#about #full_about")
const summary = document.querySelector("#about #summary")

function showAbout() {
	summary.classList.toggle("hidden")
	fullAbout.classList.toggle("hidden")
}

[fullAbout, summary].forEach(item => item.onclick = showAbout)

document.querySelectorAll("#lyrics a").forEach(item => {
	item.addEventListener("click", getAnnotation)
})

//const annotationContainer = document.getElementById("annotation-container")
/**
annotationContainer.onclick = (e) => {
	const isVisible = e.target.style.display !== "none"
	if (isVisible) {
		e.target.style.display = "none"
	}
}
**/

const card = document.getElementById("annotation")

async function getAnnotation(e) {
	e.preventDefault()
	const path = e.target.parentElement.getAttribute("href")
	const annotationId = path.match(/\/(\d+)\/.*$/)[1]
	console.log(e.offsetLeft)

	try {
		const res = await fetch(`/annotation/${annotationId}`)
		const data = await res.json()
		card.style.top = `${e.offsetY}px`
		card.style.left = `${e.offsetX}px`
		card.style.display = "block"
		card.innerHTML = data.Annotations[0].Body.HTML
	} catch (e) {
		console.error(e)
	}
}
