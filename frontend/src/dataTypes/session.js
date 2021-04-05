export const SESSION_TYPE_OPEN = 0;
export const SESSION_TYPE_PRIVATE = 1;

export default class Session {
  constructor(
    id = null,
    title = null,
    description = null,
    start = null, // ISO string
    end = null, // ISO string
    maxRiders = null,
    type = null
  ) {
    this.id = id;
    this.title = title;
    this.description = description;
    this.start = start;
    this.end = end;
    this.maximumRiders = maxRiders;
    this.type = type;
  }

  addRiders = (riders) => {
    this.riders = riders;
  };
}
