declare function examplePrintString(msg: string): Promise<void>;

declare function examplePrintStruct(input: { name: string }): Promise<void>;

declare function exampleReturnError(): Promise<void>;

declare function exampleReturnString(): Promise<string>;

declare function exampleReturnStruct(): Promise<{ name: string }>;

declare function exampleReturnStringAndError(): Promise<string>;

declare function exampleReturnStringAndNil(): Promise<string>;

declare function exampleNotification(): Promise<string>;
