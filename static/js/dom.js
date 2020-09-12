export default function CreateElement(tagName, attributes = {}, textContent = undefined) {
    const elem = document.createElement(tagName);

    for (const [key, value] of Object.entries(attributes)) {
        elem.setAttribute(key, value);
    }

    if (textContent !== undefined) {
        elem.textContent = textContent;
    }

    return elem;
}
