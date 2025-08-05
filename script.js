async function main() {
    let counter = 0;
    const numOperations = 1000000;
    
    const increment = () => {
        return new Promise(resolve => {
            setImmediate(() => {
                counter++;
                resolve();
            });
        });
    };
    
    const promises = [];
    for (let i = 0; i < numOperations; i++) {
        promises.push(increment());
    }
    
    await Promise.all(promises);
    console.log(`Final counter: ${counter}`);
}

main().catch(console.error);