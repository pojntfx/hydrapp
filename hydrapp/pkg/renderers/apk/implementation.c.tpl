#include "hydrapp_android.h"
#include <stdlib.h>
#include <jni.h>
#include <string.h>

jstring get_java_string(JNIEnv *env, char *c_string) {
  return (*env)->NewStringUTF(env, c_string);
}

char* get_c_string(JNIEnv *env, jstring java_string) {
    const char *raw_utf8 = (*env)->GetStringUTFChars(env, java_string, NULL);
    jsize java_string_len = (*env)->GetStringUTFLength(env, java_string);
    
    char *c_string = (char *)malloc(java_string_len + 1); // Add additional byte for null terminator
    if (c_string == NULL) {
        (*env)->ReleaseStringUTFChars(env, java_string, raw_utf8);

        return NULL;
    }
    
    memcpy(c_string, raw_utf8, java_string_len);
    c_string[java_string_len] = '\0'; // Add null terminator

    (*env)->ReleaseStringUTFChars(env, java_string, raw_utf8);

    return c_string;
}
