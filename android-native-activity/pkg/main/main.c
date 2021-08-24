#include "main.h"

void Java_com_example_app_MainActivity_GoLoop(JNIEnv *env, jobject activity) {
  // Start the main loop
  GoLoop(env, activity);
}

void show_toast(JNIEnv *env, jobject activity) {
  // Get the toast functions
  jclass toast_cls = (*env)->FindClass(env, "android/widget/Toast");
  jmethodID toast_cls_make_text =
      (*env)->GetStaticMethodID(env, toast_cls, "makeText",
                                "(Landroid/content/Context;Ljava/lang/"
                                "CharSequence;I)Landroid/widget/Toast;");

  // Get the message to display
  jobject msg_obj = (*env)->NewStringUTF(env, "Hello from C!");

  // Create the toast
  jobject toast_obj = (*env)->CallStaticObjectMethod(
      env, toast_cls, toast_cls_make_text, activity, msg_obj, 0);

  // Show the toast
  jclass toast_obj_class = (*env)->GetObjectClass(env, toast_obj);
  jmethodID toast_obj_class_show =
      (*env)->GetMethodID(env, toast_obj_class, "show", "()V");
  (*env)->CallVoidMethod(env, toast_obj, toast_obj_class_show);
}