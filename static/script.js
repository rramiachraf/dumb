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
	const uri = e.target.parentElement.getAttribute("href")

	xhr = new XMLHttpRequest()
	xhr.open("GET", uri + "/annotations")
	xhr.send()
	xhr.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			const parsedReponse = JSON.parse(this.responseText)
			document.getElementById("annotations").innerHTML = parsedReponse.html
		}
	}
}
