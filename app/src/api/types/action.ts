export type Action = {
  name: string;
  type: ActionType;
  url: string;
};

export type ActionType = 'navigate-to' | 'rent';
