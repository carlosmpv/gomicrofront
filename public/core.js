window.component_list = [];

customElements.define('core-plugin-slot', class extends HTMLElement {
    connectedCallback() {
        const name = this.getAttribute('name');
        const reference = `/plugins/${name}`

        fetch(reference)
            .then(r => r.text())
            .then(script => eval(script))
            .then(() => this.innerHTML = `<plugin-${name} />`)
            .catch(console.log);
    }
});