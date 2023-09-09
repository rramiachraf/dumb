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
	console.log("Annotations are not yet implemented!", uri)

	xhr = new XMLHttpRequest()
	xhr.open("GET", uri + "/annotations")
	xhr.send()
	xhr.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			json = JSON.parse(this.responseText)
			alert(json.html)
			// TODO: display annotations properly
		}
	}
}
