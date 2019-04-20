(() => {
    function main() {
        const reactRoot = document.getElementById('react-root');
        const one = document.createElement('div');
        one.className = 'one';
        reactRoot.appendChild(one);

        const stiffness = 0.0002;
        const dampingRatio = 0.5;
        const dampingCoefficient = dampingRatio * 2 * Math.sqrt(stiffness);

        let dragX = 0;
        let dragY = 0;
        let dragVelocityX = 0;
        let dragVelocityY = 0;
        let dragPreviousTime = 0;

        let simulation = null;

        const drag = new Drag(
            one,
            () => {
                if (simulation !== null) {
                    simulation.stop();
                }
            },
            (x, y) => {
                one.style.transform = 'translate(' + x + 'px, ' + y + 'px)';

                const now = performance.now();
                const dt = now - dragPreviousTime;
                dragVelocityX = (x - dragX) / dt;
                dragVelocityY = (y - dragY) / dt;
                dragX = x;
                dragY = y;
                dragPreviousTime = now;
            },
            () => {
                const springX = new Spring(stiffness, dampingCoefficient, dragVelocityX, dragX);
                const springY = new Spring(stiffness, dampingCoefficient, dragVelocityY, dragY);

                simulation = new Simulation(1000 / 60, dt => {
                    springX.step(dt);
                    springY.step(dt);

                    one.style.transform = 'translate(' + springX.displacement + 'px, ' + springY.displacement + 'px)';

                    const xStopped = Math.abs(springX.velocity) < 0.001 && Math.abs(springX.acceleration) < 0.001;
                    const yStopped = Math.abs(springY.velocity) < 0.001 && Math.abs(springY.acceleration) < 0.001;

                    return !(xStopped && yStopped);
                });
            },
        );
    }

    main();
})();
