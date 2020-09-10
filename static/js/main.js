// TODO: Make HMS ships stand out more, make other vessels take up less screen space.

const apiRoute = '/api/movements';

function fromToSentence(day) {
    if (day.from.name !== '' && day.to.name !== '') {
        return `${day.from.name} to ${day.to.name}`;
    }
    return '';
}

function createCard(day) {
    const card = document.createElement('div');
    card.setAttribute('class', 'card horizontal');

    const cardImage = document.createElement('div');
    const cardImageImg = document.createElement('img');
    cardImage.setAttribute('class', 'card-image');
    cardImageImg.setAttribute('src', day.imageUrl);
    cardImage.appendChild(cardImageImg);

    const cardStacked = document.createElement('div');
    const cardContent = document.createElement('div');
    cardStacked.setAttribute('class', 'card-stacked');
    cardContent.setAttribute('class', 'card-content');

    const titleH6 = document.createElement('h6');
    const timeP = document.createElement('p');
    const moveP = document.createElement('p');
    titleH6.textContent = day.name;
    timeP.textContent = day.time;
    moveP.textContent = fromToSentence(day);

    cardContent.appendChild(titleH6);
    cardContent.appendChild(timeP);
    cardContent.appendChild(moveP);
    cardStacked.appendChild(cardContent);

    card.appendChild(cardImage);
    card.appendChild(cardStacked);
    return card;
}

function addRow(today = undefined, tomorrow = undefined) {
    const container = document.querySelector('.container');
    const row = document.createElement('div');
    row.setAttribute('class', 'row');

    if (today !== undefined) {
        const colL = document.createElement('div');
        colL.setAttribute('class', 'col s6');
        const cardL = createCard(today);
        colL.appendChild(cardL);
        row.appendChild(colL);
    }

    if (tomorrow !== undefined) {
        const colR = document.createElement('div');
        if (today === undefined) {
            colR.setAttribute('class', 'col s6 offset-s6');
        } else {
            colR.setAttribute('class', 'col s6');
        }
        const cardR = createCard(tomorrow);
        colR.appendChild(cardR);
        row.appendChild(colR);
    }

    container.appendChild(row);
}

document.addEventListener('DOMContentLoaded', async () => {
    const r = await fetch(apiRoute);
    const movements = await r.json();

    let i = 0;
    for (; i < movements.today.length; i++) {
        const todayMovement = movements.today[i];
        const tomorrowMovement = movements.tomorrow[i];
        addRow(todayMovement, tomorrowMovement);
    }

    if (movements.tomorrow.length > movements.today.length) {
        for (; i < movements.tomorrow.length; i++) {
            const tomorrowMovement = movements.tomorrow[i];
            addRow(undefined, tomorrowMovement);
        }
    }
});
