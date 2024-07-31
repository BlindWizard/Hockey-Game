import Cookie from 'js-cookie';

export default class PlayerId {
    getId() {
        return Cookie.get('player-id') || null;
    }

    storeId(id) {
        Cookie.set('player-id', id, {expires: 365});
    }
}