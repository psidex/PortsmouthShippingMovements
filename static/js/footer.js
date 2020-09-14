// Update the footer text with the current version.
export default async function UpdateFooter() {
    const footer = document.getElementById('footerText');
    const version = await (await fetch('/version')).text();
    footer.textContent = `github.com/psidex/PortsmouthShippingMovements @ ${version}`;
}
