const componentName = 'hello-world'

customElements.define(`plugin-${componentName}`, class extends HTMLElement {
    connectedCallback() {
        this.mount().then(() => this.onMount())
    }

    async mount() {
        const reference = `/plugins/${componentName}?template`
        await fetch(reference)
            .then(r => r.text())
            .then(t => this.innerHTML = t)
            .catch(console.error);
    }

    onMount() {
        const dynamicValue = this.querySelector("#dynamic-value");
        
        const interval = setInterval(() => {
            dynamicValue.innerHTML = (new Date()).toLocaleTimeString("pt-Br")
        }, 1000)
    }
});

window.component_list = [...(window.component_list || []), componentName];