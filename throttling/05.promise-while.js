
class PromiseWhile {
  constructor(context) {
    this.context = context;
  }

  // condition - function returning bool
  // action - function returning promise
  Await(condition, action) {
    return new Promise((resolve, reject) => {
      let returning;
      const loop = () => {
        if (!condition.call(this.context)) {
          return resolve(returning);
        }
        return action.call(this.context)
          .then((data) => {
            returning = data;
            setTimeout(() => loop(), this.context.timeoutBetweenCalls !== undefined
              ? this.context.timeoutBetweenCalls
              : 1000);
          })
          .catch(err =>
            reject(err));
      };

      process.nextTick(loop);
    });
  }
}

module.exports = PromiseWhile;