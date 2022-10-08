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

function getAnnotation(e) {
	e.preventDefault()
	console.log(e.target.parentElement.getAttribute("href"))
}
