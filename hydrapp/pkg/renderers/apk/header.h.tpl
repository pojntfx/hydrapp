#include <jni.h>

jstring get_java_string(JNIEnv *env, char *c_string);

char* get_c_string(JNIEnv *env, jstring java_string);
