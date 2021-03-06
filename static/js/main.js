import CreateElement from './dom.js';
import UpdateFooter from './footer.js';

const apiRoute = '/api/movements';

function fromToSentence(movement) {
    if (movement.type === 0) {
        return `Moving from ${movement.from.name} to ${movement.to.name}`;
    }
    return '';
}

function fromToAbbreviation(movement) {
    if (movement.type === 0) {
        return `${movement.from.abbreviation} to ${movement.to.abbreviation}`;
    }
    return '';
}

// Constructs a <li> element for a movement.
function addLi(movement, collapsible) {
    const movementLi = CreateElement('li');
    if (movement.name.startsWith('HMS')) {
        // This ship is Navy so make it stand out.
        movementLi.setAttribute('class', 'active');
    }

    //
    // HEADER
    //
    const headerDiv = CreateElement('div', { class: 'collapsible-header' });
    const movementTitleP = CreateElement('p', { class: 'movement-title' }, `${movement.time} - ${movement.name}`);
    const fromToP = CreateElement('p', { class: 'from-to-abbrv' }, fromToAbbreviation(movement));

    if (!movement.name.startsWith('MV')) {
        // All the ferries are "MV", anything not a ferry could be interesting.
        movementTitleP.setAttribute('class', 'movement-title text-bold');
    }

    headerDiv.appendChild(movementTitleP);
    headerDiv.appendChild(fromToP);

    movementLi.appendChild(headerDiv);

    // If it's a ship movement and not a notice, add a body.
    if (movement.type === 0) {
        //
        // BODY
        //
        const bodyDiv = CreateElement('div', { class: 'collapsible-body' });
        const bodyImg = CreateElement('img', { src: movement.imageUrl });
        const bodyP = CreateElement('p', {}, fromToSentence(movement));
        const vesselFinderA = CreateElement('a', {
            href: movement.infoUrl,
            target: '_blank',
            class: 'tooltipped',
            'data-position': 'bottom',
            'data-tooltip': 'Search on Vessel Finder',
        });
        const vesselFinderImg = CreateElement('img', {
            class: 'info-link-img',
            src: '/images/compass.svg',
        });

        vesselFinderA.appendChild(vesselFinderImg);

        bodyDiv.appendChild(bodyImg);
        bodyDiv.appendChild(bodyP);
        bodyDiv.appendChild(vesselFinderA);

        movementLi.appendChild(bodyDiv);
    }

    collapsible.appendChild(movementLi);
}

document.addEventListener('DOMContentLoaded', async () => {
    await UpdateFooter();

    const todayCollapsible = document.querySelector('#todayCollapsible');
    const tomorrowCollapsible = document.querySelector('#tomorrowCollapsible');

    const r = await fetch(apiRoute);
    const movements = await r.json();

    if (movements.today === null) {
        movements.today = [];
    }
    if (movements.tomorrow === null) {
        movements.tomorrow = [];
    }

    let i = 0;
    for (; i < movements.today.length; i++) {
        const todayMovement = movements.today[i];
        const tomorrowMovement = movements.tomorrow[i];
        addLi(todayMovement, todayCollapsible);
        if (tomorrowMovement !== undefined) {
            addLi(tomorrowMovement, tomorrowCollapsible);
        }
    }

    if (movements.tomorrow.length > movements.today.length) {
        for (; i < movements.tomorrow.length; i++) {
            const tomorrowMovement = movements.tomorrow[i];
            addLi(tomorrowMovement, tomorrowCollapsible);
        }
    }

    // Init after all the HTML is ready.
    const collapsibleElems = document.querySelectorAll('.collapsible');
    const tooltippedElems = document.querySelectorAll('.tooltipped');
    M.Collapsible.init(collapsibleElems, { accordion: false });
    M.Tooltip.init(tooltippedElems, {});
});
