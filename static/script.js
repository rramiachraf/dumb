document.querySelectorAll("#lyrics a").forEach(item => {
	item.addEventListener("click", getAnnotation)
})

function getAnnotation(e) {
	e.preventDefault()
}
