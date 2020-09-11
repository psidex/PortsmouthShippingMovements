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
    const movementLi = document.createElement('li');

    //
    // HEADER
    //
    const headerDiv = document.createElement('div');
    const movementTitleP = document.createElement('p');
    const fromToP = document.createElement('p');

    headerDiv.setAttribute('class', 'collapsible-header');

    movementTitleP.textContent = `${movement.time} - ${movement.name}`;
    movementTitleP.setAttribute('class', 'movement-title');
    if (!movement.name.startsWith('MV')) {
        // All the ferries are "MV", anything not a ferry could be interesting.
        movementTitleP.setAttribute('class', 'movement-title text-bold');
    }

    fromToP.textContent = fromToAbbreviation(movement);
    fromToP.setAttribute('class', 'right-margin');

    headerDiv.appendChild(movementTitleP);
    headerDiv.appendChild(fromToP);
    movementLi.appendChild(headerDiv);

    // If it's a ship movement and not a notice, add a body.
    if (movement.type === 0) {
        //
        // BODY
        //
        const bodyDiv = document.createElement('div');
        const bodyImg = document.createElement('img');
        const bodyP = document.createElement('p');
        const bodyInfoA = document.createElement('a');
        const bodyInfoImg = document.createElement('img');

        bodyDiv.setAttribute('class', 'collapsible-body');
        bodyImg.setAttribute('src', movement.imageUrl);
        bodyP.textContent = fromToSentence(movement);
        bodyInfoA.setAttribute('href', movement.infoUrl);
        bodyInfoA.setAttribute('target', '_blank');
        bodyInfoA.setAttribute('class', 'tooltipped');
        bodyInfoA.setAttribute('data-position', 'bottom');
        bodyInfoA.setAttribute('data-tooltip', 'Vessel Finder');
        bodyInfoImg.setAttribute('class', 'info-link-img');
        bodyInfoImg.setAttribute('src', '/images/compass.svg');

        bodyInfoA.appendChild(bodyInfoImg);

        bodyDiv.appendChild(bodyImg);
        bodyDiv.appendChild(bodyP);
        bodyDiv.appendChild(bodyInfoA);
        movementLi.appendChild(bodyDiv);
    }

    if (movement.name.startsWith('HMS')) {
        movementLi.setAttribute('class', 'active');
    }
    collapsible.appendChild(movementLi);
}

document.addEventListener('DOMContentLoaded', async () => {
    const todayCollapsible = document.querySelector('#todayCollapsible');
    const tomorrowCollapsible = document.querySelector('#tomorrowCollapsible');

    const r = await fetch(apiRoute);
    const movements = await r.json();

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

    // Init after all the HTML is setup.
    const collapsibleElems = document.querySelectorAll('.collapsible');
    M.Collapsible.init(collapsibleElems, {accordion: false});

    const tooltippedElems = document.querySelectorAll('.tooltipped');
    M.Tooltip.init(tooltippedElems, {});
});
