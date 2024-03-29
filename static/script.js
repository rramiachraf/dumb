const fullAbout = document.querySelector("#about #full_about")
const summary = document.querySelector("#about #summary")

function showAbout() {
	summary.classList.toggle("hidden")
	fullAbout.classList.toggle("hidden")
}

fullAbout && [fullAbout, summary].forEach(item => item.onclick = showAbout)

document.querySelectorAll("#lyrics a").forEach(item => {
	item.addEventListener("click", getAnnotation)
})

function getAnnotation(e) {
	e.preventDefault()
	//document.querySelector('.annotation')?.remove()
	const uri = e.target.parentElement.getAttribute("href")
	const presentAnnotation = document.getElementById(uri)
	if (presentAnnotation) {
		presentAnnotation.remove()
		return
	}

	xhr = new XMLHttpRequest()
	xhr.open("GET", uri + "/annotations")
	xhr.send()
	xhr.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			const parsedReponse = JSON.parse(this.responseText)
			const annotationDiv = document.createElement('div');
			annotationDiv.innerHTML = parsedReponse.html
			annotationDiv.id = uri
			annotationDiv.className = "annotation"
			e.target.parentElement.insertAdjacentElement('afterend', annotationDiv)
		}
	}
}

window._currentTheme = localStorage.getItem("_theme") || "light"
setTheme(window._currentTheme)

const themeChooser = document.getElementById("choose-theme")
themeChooser.addEventListener("click", function() {
	if (window._currentTheme === "dark") {
		setTheme("light")
	} else {
		setTheme("dark")
	}
})

function setTheme(theme) {
	const toggler = document.getElementById("ic_fluent_dark_theme_24_regular")
	switch (theme) {
		case "dark":
			toggler.setAttribute("fill", "#fff")
			localStorage.setItem("_theme", "dark")
			document.body.classList.add("dark")
			window._currentTheme = "dark"
			return
		case "light":
			toggler.setAttribute("fill", "#181d31")
			localStorage.setItem("_theme", "light")
			document.body.classList.remove("dark")
			window._currentTheme = "light"
			return

	}
}
