import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    scenarios: {
        reactions: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '5s', target: 100 },
                { duration: '5s', target: 0 },
            ],
            gracefulRampDown: '0s',
            exec: 'reactions',
        },
    },
};


export function reactions() {
    sleep(0.5 + Math.random());

    for (let i = 1; i <= 5; i++) {
        const availableReactions = http.get(`http://api:8080/livestreams/${i}/reactions`);
        check(availableReactions, {
            'status is 200': (r) => r.status === 200,
            'at least one reaction available': (r) => r.json().length > 0,
        });

        const reaction = availableReactions.json()[Math.floor(Math.random() * availableReactions.json().length)];

        const res = http.post(`http://api:8080/livestreams/${i}/reactions/${reaction}`);
        check(res, {'status is 200': (r) => r.status === 200});
    }
}
